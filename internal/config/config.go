package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type (
	Config struct {
		Web  *Web
		DB   *DB
		Stan *Stan
	}

	Web struct {
		Host          string
		Port          int32
		WaitShutdown  time.Duration
		AllowedOrigin string
	}

	DB struct {
		User         string
		Password     string
		Host         string
		Port         int
		DatabaseName string
	}
	Stan struct {
		URL       string
		ClusterID string
		ClientID  string
		Subject   string
	}
)

func New() *Config {
	return &Config{
		Web: &Web{
			Host: "0.0.0.0",
			Port: 8080,
		},
		DB: &DB{
			User:         GetEnv("POSTGRES_USER", "postgres"),
			Password:     GetEnv("POSTGRES_PASSWORD", "23785"),
			Host:         GetEnv("DB_HOST", "0.0.0.0"),
			Port:         GetnvInt("DB_PORT", 5432),
			DatabaseName: GetEnv("POSTGRES_DB", "10"),
		},
		Stan: &Stan{
			URL:       GetEnv("STAN_SERVER", "http://127.0.0.1:4222"),
			ClusterID: GetEnv("CLUSTER_ID", "test-cluster"),
			ClientID:  GetEnv("CLIENT_ID", "stan-sub"),
			Subject:   GetEnv("SUBJECT", "orders"),
		},
	}
}

// (w *Web) если передается по указателю, то менять можно
// (w Web) если передается по значению, то меняется копия и её нужно будет куда-то сохранять
func (w Web) Address() string {
	return fmt.Sprintf("%s:%d", w.Host, w.Port)
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
func GetnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("$%s not parse as int.\n", key)
		return defaultValue
	}
	return result
}
