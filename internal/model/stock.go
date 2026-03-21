package model

import "time"

type Stock struct {
	ID          int    `json:"id"`
	Symbol      string `json:"symbol"`
	LastClose   string `json:"last_close"`
	High52W     string `json:"high_52w"`
	Low52w      string `json:"low_52w"`
	DiffPercent string `json:"diff"`
	CurDate     string `json:"cur_date"`

	// 仅为兼容当前 repo 模板代码
	Name      string    `json:"name,omitempty"`
	Price     float64   `json:"price,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
