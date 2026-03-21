package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"project/internal/api"
	"project/internal/api/middleware"
	"project/internal/config"
	"project/internal/db"
	"project/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warn: .env not loaded: %v", err)
	}

	cfg := config.Load()
	app := fiber.New()

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

	app.Use(middleware.Logger(l))
	api.RegisterRoutes(app, l, mysqlDB, cfg)
	l.Info("server start", zap.String("port", cfg.Port))
	log.Fatal(app.Listen(":" + cfg.Port))
}
