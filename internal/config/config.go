package config

import "os"

type Config struct {
	Port           string
	LogLevel       string
	AppName        string
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
	DBName         string
	APIKey         string
	AlphaBaseURL   string
	SyncSymbols    string
	SyncCronExpr   string
	SyncRunOnStart bool
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	v := getEnv(key, "")
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "TRUE", "True", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "off", "OFF":
		return false
	default:
		return def
	}
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		AppName:        getEnv("APP_NAME", "stock-app"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBName:         getEnv("DB_NAME", "stock_app"),
		APIKey:         getEnv("ALPHA_API_KEY", "test_api_key"),
		AlphaBaseURL:   getEnv("ALPHA_BASE_URL", "https://www.alphavantage.co/query"),
		SyncSymbols:    getEnv("SYNC_SYMBOLS", "MSFT,AAPL,META"),
		SyncCronExpr:   getEnv("SYNC_CRON", "*/30 * * * *"),
		SyncRunOnStart: getEnvBool("SYNC_RUN_ON_START", true),
	}
}
