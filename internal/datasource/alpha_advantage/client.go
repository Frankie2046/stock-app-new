package alpha_advantage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(httpClient *http.Client, baseURL, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://www.alphavantage.co/query"
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (c *Client) GetWeeklyAdjusted(ctx context.Context, symbol string) ([]WeeklyBar, string, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, "", fmt.Errorf("symbol is empty")
	}
	if strings.TrimSpace(c.apiKey) == "" {
		return nil, "", fmt.Errorf("ALPHA_API_KEY is empty")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, "", fmt.Errorf("invalid base url: %w", err)
	}
	q := u.Query()
	q.Set("function", "TIME_SERIES_WEEKLY_ADJUSTED")
	q.Set("symbol", symbol)
	q.Set("apikey", c.apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("alphavantage status=%d body=%s", resp.StatusCode, string(body))
	}

	var parsed weeklyAdjustedResp
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, "", err
	}
	if parsed.ErrorMessage != "" {
		return nil, "", fmt.Errorf("alphavantage error: %s", parsed.ErrorMessage)
	}
	if parsed.Note != "" {
		return nil, "", fmt.Errorf("alphavantage note: %s", parsed.Note)
	}
	if len(parsed.Weekly) == 0 {
		return nil, "", fmt.Errorf("weekly time series is empty")
	}

	bars := make([]WeeklyBar, 0, len(parsed.Weekly))
	for date, row := range parsed.Weekly {
		high, err := strconv.ParseFloat(row.High, 64)
		if err != nil {
			return nil, "", fmt.Errorf("parse high failed on %s: %w", date, err)
		}
		low, err := strconv.ParseFloat(row.Low, 64)
		if err != nil {
			return nil, "", fmt.Errorf("parse low failed on %s: %w", date, err)
		}
		adj, err := strconv.ParseFloat(row.AdjustedClose, 64)
		if err != nil {
			return nil, "", fmt.Errorf("parse adjusted close failed on %s: %w", date, err)
		}
		bars = append(bars, WeeklyBar{
			Date:          date,
			High:          high,
			Low:           low,
			AdjustedClose: adj,
		})
	}

	// yyyy-mm-dd 可以直接字符串倒序
	sort.Slice(bars, func(i, j int) bool {
		return bars[i].Date > bars[j].Date
	})

	return bars, parsed.MetaData.LastRefreshed, nil
}
