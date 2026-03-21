package model

type Stock struct {
	ID          int64   `json:"id"`
	Symbol      string  `json:"symbol"`
	LastClose   float64 `json:"last_close"`
	High52W     float64 `json:"high_52w"`
	Low52W      float64 `json:"low_52w"`
	DiffPercent float64 `json:"diff_percent"`
	CurDate     string  `json:"cur_date"`
	RefreshTime string  `json:"refresh_time"`
}
