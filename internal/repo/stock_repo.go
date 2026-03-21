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
	query := `SELECT id, symbol, last_close, high_52w, low_52w, percent_diff, cur_date FROM stocks WHERE id = ? LIMIT 1`
	var s model.Stock
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&s.ID, &s.Symbol, &s.LastClose, &s.High52W, &s.Low52w, &s.DiffPercent, &s.CurDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *StockRepo) GetBySymbol(ctx context.Context, symbol string) (*model.Stock, error) {
	query := `SELECT id, symbol, name, price, created_at, updated_at FROM stocks WHERE symbol = ?`
	var s model.Stock
	err := r.db.QueryRowContext(ctx, query, symbol).
		Scan(&s.ID, &s.Symbol, &s.Name, &s.Price, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *StockRepo) List(ctx context.Context) ([]model.Stock, error) {
	query := `SELECT id, symbol, name, price, created_at, updated_at FROM stocks ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Stock
	for rows.Next() {
		var s model.Stock
		if err := rows.Scan(&s.ID, &s.Symbol, &s.Name, &s.Price, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *StockRepo) UpdatePrice(ctx context.Context, symbol string, price float64) error {
	query := `UPDATE stocks SET price = ? WHERE symbol = ?`
	res, err := r.db.ExecContext(ctx, query, price, symbol)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *StockRepo) DeleteBySymbol(ctx context.Context, symbol string) error {
	query := `DELETE FROM stocks WHERE symbol = ?`
	res, err := r.db.ExecContext(ctx, query, symbol)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
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
		s.Low52w,
		s.DiffPercent,
		"AlphaVantage",
		s.CurDate,
	)
	return err
}
