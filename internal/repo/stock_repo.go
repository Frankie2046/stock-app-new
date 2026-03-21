package repo

import (
	"context"
	"database/sql"
	"errors"

	"project/internal/model"
)

var ErrNotFound = errors.New("stock not found")

type StockRepo struct {
	db *sql.DB
}

func NewStockRepo(db *sql.DB) *StockRepo {
	return &StockRepo{db: db}
}

func (r *StockRepo) GetByID(ctx context.Context, id int) (*model.Stock, error) {
	query := `SELECT id, symbol, last_close, high_52w, low_52w, percent_diff, DATE_FORMAT(cur_date, '%Y-%m-%d'), DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') FROM stocks WHERE id = ? LIMIT 1`
	var s model.Stock
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&s.ID, &s.Symbol, &s.LastClose, &s.High52W, &s.Low52W, &s.DiffPercent, &s.CurDate, &s.RefreshTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *StockRepo) GetBySymbol(ctx context.Context, symbol string) (*model.Stock, error) {
	query := `SELECT id, symbol, last_close, high_52w, low_52w, percent_diff, DATE_FORMAT(cur_date, '%Y-%m-%d'), DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') FROM stocks WHERE symbol = ? ORDER BY cur_date DESC LIMIT 1`
	var s model.Stock
	err := r.db.QueryRowContext(ctx, query, symbol).
		Scan(&s.ID, &s.Symbol, &s.LastClose, &s.High52W, &s.Low52W, &s.DiffPercent, &s.CurDate, &s.RefreshTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *StockRepo) ListLatestDateStocks(ctx context.Context) ([]model.Stock, error) {
	query := `
	SELECT s.id, s.symbol, s.last_close, s.high_52w, s.low_52w, s.percent_diff,
		DATE_FORMAT(s.cur_date, '%Y-%m-%d') AS cur_date,
		DATE_FORMAT(s.updated_at, '%Y-%m-%d %H:%i:%s') AS refresh_time
	FROM stocks s
	JOIN (
	SELECT symbol, MAX(cur_date) AS max_date
	FROM stocks
	GROUP BY symbol
	) t
	ON s.symbol = t.symbol
	AND s.cur_date = t.max_date
	ORDER BY s.symbol,s.updated_at DESC;
`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.Stock, 0)
	for rows.Next() {
		var s model.Stock
		if err := rows.Scan(&s.ID, &s.Symbol, &s.LastClose, &s.High52W, &s.Low52W, &s.DiffPercent, &s.CurDate, &s.RefreshTime); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *StockRepo) UpsertDaily(ctx context.Context, s model.Stock) error {
	query := `
INSERT INTO stocks (
  symbol, last_close, high_52w, low_52w, percent_diff, is_active, notes, cur_date
) VALUES (?, ?, ?, ?, ?, 1, ?, ?)
ON DUPLICATE KEY UPDATE
  last_close = VALUES(last_close),
  high_52w = VALUES(high_52w),
  low_52w = VALUES(low_52w),
  percent_diff = VALUES(percent_diff),
  notes = VALUES(notes),
  updated_at = CURRENT_TIMESTAMP
`
	_, err := r.db.ExecContext(
		ctx,
		query,
		s.Symbol,
		s.LastClose,
		s.High52W,
		s.Low52W,
		s.DiffPercent,
		"AlphaVantage",
		s.CurDate,
	)
	return err
}
