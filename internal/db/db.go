package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Params   string // 例如: charset=utf8mb4&parseTime=True&loc=Local
}

func NewMySQL(c Config) (*sql.DB, error) {
	if c.Port == "" {
		c.Port = "3306"
	}
	if c.Params == "" {
		c.Params = "charset=utf8mb4&parseTime=True&loc=Local"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Params)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 连接池参数可按压测调整
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
