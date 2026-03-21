package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"project/internal/datasource/alpha_advantage"
	"project/internal/model"
	"project/internal/repo"
)

type MarketDataProvider interface {
	GetWeeklyAdjusted(ctx context.Context, symbol string) ([]alpha_advantage.WeeklyBar, string, error)
}

type StockService struct {
	repo     *repo.StockRepo
	provider MarketDataProvider
}

func NewStockService(repo *repo.StockRepo, provider MarketDataProvider) *StockService {
	return &StockService{
		repo:     repo,
		provider: provider,
	}
}

func (s *StockService) GetStock(id int) (*model.Stock, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stock, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, errors.New("stock not found")
		}
		return nil, err
	}

	return stock, nil
}

func (s *StockService) GetLatestDateStocks() ([]model.Stock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.repo.ListLatestDateStocks(ctx)
}

func (s *StockService) SyncSymbol(symbol string) error {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return errors.New("symbol is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bars, lastRefreshed, err := s.provider.GetWeeklyAdjusted(ctx, symbol)
	if err != nil {
		return err
	}
	if len(bars) == 0 {
		return errors.New("no weekly data")
	}

	n := 52
	if len(bars) < n {
		n = len(bars)
	}

	lastClose := bars[0].AdjustedClose
	high52 := bars[0].High
	low52 := bars[0].Low

	for i := 1; i < n; i++ {
		if bars[i].High > high52 {
			high52 = bars[i].High
		}
		if bars[i].Low < low52 {
			low52 = bars[i].Low
		}
	}
	if high52 <= 0 {
		return errors.New("invalid high_52w <= 0")
	}

	diff := (lastClose - high52) / high52 * 100

	curDate := strings.TrimSpace(lastRefreshed)
	if len(curDate) >= 10 {
		curDate = curDate[:10]
	} else {
		curDate = bars[0].Date
	}

	stock := model.Stock{
		Symbol:      symbol,
		LastClose:   lastClose,
		High52W:     high52,
		Low52W:      low52,
		DiffPercent: diff,
		CurDate:     curDate,
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCancel()

	return s.repo.UpsertDaily(dbCtx, stock)
}
