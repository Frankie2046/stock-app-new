package alpha_advantage

type WeeklyBar struct {
	Date          string
	High          float64
	Low           float64
	AdjustedClose float64
}

type weeklyAdjustedResp struct {
	MetaData struct {
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
	} `json:"Meta Data"`
	Weekly map[string]struct {
		High          string `json:"2. high"`
		Low           string `json:"3. low"`
		AdjustedClose string `json:"5. adjusted close"`
	} `json:"Weekly Adjusted Time Series"`
	Note         string `json:"Note"`
	ErrorMessage string `json:"Error Message"`
}
