package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"project/internal/api/handler"
	"project/internal/repo"
	"project/internal/service"
)

func RegisterRoutes(app *fiber.App, logger *zap.Logger, mysqlDB *sql.DB) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	stockRepo := repo.NewStockRepo(mysqlDB)
	stockService := service.NewStockService(stockRepo)
	stockHandler := handler.NewStockHandler(stockService, logger)

	app.Get("/stocks/:id", stockHandler.GetStock)
}
