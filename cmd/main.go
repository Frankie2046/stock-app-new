package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	template "github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"project/internal/api"
	"project/internal/api/middleware"
	"project/internal/config"
	"project/internal/datasource/alpha_advantage"
	"project/internal/db"
	"project/internal/repo"
	"project/internal/scheduler"
	"project/internal/service"
	"project/pkg/logger"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warn: .env not loaded: %v", err)
	}

	cfg := config.Load()
	engine := template.New("./templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	l, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer func() { _ = l.Sync() }()

	mysqlDB, err := db.NewMySQL(db.Config{
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Name:     cfg.DBName,
	})
	if err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}
	defer mysqlDB.Close()

	stockRepo := repo.NewStockRepo(mysqlDB)
	alphaClient := alpha_advantage.NewClient(&http.Client{}, cfg.AlphaBaseURL, cfg.APIKey)
	stockService := service.NewStockService(stockRepo, alphaClient)

	cronRunner, err := scheduler.StartStockSync(cfg, l, stockService)
	if err != nil {
		log.Fatalf("start stock sync scheduler failed: %v", err)
	}
	if cronRunner != nil {
		defer cronRunner.Stop()
	}

	app.Use(middleware.Logger(l))
	api.RegisterRoutes(app, l, stockService)

	l.Info("server start", zap.String("port", cfg.Port))
	log.Fatal(app.Listen(":" + cfg.Port))
}
