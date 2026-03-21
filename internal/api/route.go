package api

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"project/internal/api/handler"
	"project/internal/service"
)

func RegisterRoutes(app *fiber.App, logger *zap.Logger, stockService *service.StockService) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	stockHandler := handler.NewStockHandler(stockService, logger)

	app.Get("/", stockHandler.Dashboard)
	app.Get("/stocks/:id", stockHandler.GetStock)
	app.Post("/stocks/sync/:symbol", stockHandler.SyncSymbol)
}
