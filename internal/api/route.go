package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"project/internal/api/handler"
	"project/internal/config"
	"project/internal/datasource/alpha_advantage"
	"project/internal/repo"
	"project/internal/service"
)

func RegisterRoutes(app *fiber.App, logger *zap.Logger, mysqlDB *sql.DB, cfg *config.Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	stockRepo := repo.NewStockRepo(mysqlDB)
	alphaClient := alpha_advantage.NewClient(&http.Client{Timeout: 15 * time.Second}, cfg.AlphaBaseURL, cfg.APIKey)
	stockService := service.NewStockService(stockRepo, alphaClient)
	stockHandler := handler.NewStockHandler(stockService, logger)

	app.Get("/", stockHandler.Dashboard)
	app.Get("/stocks/:id", stockHandler.GetStock)
	app.Post("/stocks/sync/:symbol", stockHandler.SyncSymbol)
}
