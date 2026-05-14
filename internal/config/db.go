package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	host     string
	user     string
	password string
	dbname   string
	port     string
	sslmode  string
}

func (c DBConfig) dsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.host, c.user, c.password, c.dbname, c.port, c.sslmode)
}

func (c DBConfig) OpenConnection() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		host:     getenv("DB_HOST", "localhost"),
		user:     getenv("DB_USER", "postgres"),
		password: getenv("DB_PASSWORD", "postgres"),
		dbname:   getenv("DB_NAME", "test_db"),
		port:     getenv("DB_PORT", "5432"),
		sslmode:  getenv("DB_SSLMODE", "disable"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
