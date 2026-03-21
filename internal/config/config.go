package config

import "os"

type Config struct {
	Port       string
	LogLevel   string
	AppName    string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		AppName:    getEnv("APP_NAME", "stock-app"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "stock_app"),
	}
}
