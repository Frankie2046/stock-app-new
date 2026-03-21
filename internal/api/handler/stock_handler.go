package handler

import (
	"strconv"
	"strings"
	"time"

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
		h.logger.Warn("invalid stock id param", zap.String("id", c.Params("id")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	stock, err := h.service.GetStock(id)
	if err != nil {
		h.logger.Warn("get stock failed", zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stock)
}

func (h *StockHandler) SyncSymbol(c *fiber.Ctx) error {
	symbol := strings.ToUpper(strings.TrimSpace(c.Params("symbol")))
	if symbol == "" {
		h.logger.Warn("invalid stock symbol param", zap.String("symbol", c.Params("symbol")))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid symbol"})
	}

	if err := h.service.SyncSymbol(symbol); err != nil {
		h.logger.Warn("sync stock failed", zap.String("symbol", symbol), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "sync success", "symbol": symbol})
}

func (h *StockHandler) Dashboard(c *fiber.Ctx) error {
	stocks, err := h.service.GetLatestDateStocks()
	if err != nil {
		h.logger.Error("load dashboard failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("load dashboard failed")
	}

	pageDate := time.Now().Format("2006-01-02")
	if len(stocks) > 0 && stocks[0].CurDate != "" {
		pageDate = stocks[0].CurDate
	}

	return c.Render("index", fiber.Map{
		"Stocks": stocks,
		"Date":   pageDate,
	})
}
