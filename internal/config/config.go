package config

import (
	"fmt"
	"os"
)

type Config struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DBName  string
	SSLMode string
}

func Load() Config {
	return Config{
		Host:    getenv("PG_HOST", "localhost"),
		Port:    getenv("PG_PORT", "5430"),
		User:    getenv("PG_USER", "user"),
		Pass:    getenv("PG_PASSWORD", "password"),
		DBName:  getenv("PG_DB", "mydatabase"),
		SSLMode: getenv("PG_SSLMODE", "disable"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Pass, c.Host, c.Port, c.DBName, c.SSLMode)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
