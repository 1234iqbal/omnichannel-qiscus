package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Port           string
	RedisURL       string
	PostgresURL    string
	RequestTimeout time.Duration
	QiscusConfig   QiscusConfig
}

type QiscusConfig struct {
	BaseURL   string
	AppID     string
	SecretKey string
	Timeout   time.Duration
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	postgresURL := os.Getenv("POSTGRES_URL")

	// Qiscus configuration
	qiscusBaseURL := os.Getenv("QISCUS_BASE_URL")
	if qiscusBaseURL == "" {
		qiscusBaseURL = "https://omnichannel.qiscus.com"
	}

	log.Println("Qiscus App ID:", os.Getenv("QISCUS_APP_ID"))
	log.Println("Qiscus Secret Key:", os.Getenv("QISCUS_SECRET_KEY"))

	return &Config{
		Port:           port,
		RedisURL:       redisURL,
		PostgresURL:    postgresURL,
		RequestTimeout: 60 * time.Second,
		QiscusConfig: QiscusConfig{
			BaseURL:   qiscusBaseURL,
			AppID:     os.Getenv("QISCUS_APP_ID"),
			SecretKey: os.Getenv("QISCUS_SECRET_KEY"),
			Timeout:   30 * time.Second,
		},
	}
}
