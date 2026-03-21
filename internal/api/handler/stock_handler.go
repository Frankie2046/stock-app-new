package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"project/internal/service"
)

type StockHandler struct {
	service *service.StockService
	logger  *zap.Logger
}

func NewStockHandler(service *service.StockService, logger *zap.Logger) *StockHandler {
	return &StockHandler{
		service: service,
		logger:  logger,
	}
}

func (h *StockHandler) GetStock(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Warn("invalid stock id param",
			zap.String("id", c.Params("id")),
		)
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	stock, err := h.service.GetStock(id)
	if err != nil {
		h.logger.Warn("get stock failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stock)
}
