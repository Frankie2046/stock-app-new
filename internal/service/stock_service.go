package service

import (
	"context"
	"errors"
	"time"

	"project/internal/model"
	"project/internal/repo"
)

type StockService struct {
	repo *repo.StockRepo
}

func NewStockService(repo *repo.StockRepo) *StockService {
	return &StockService{repo: repo}
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
