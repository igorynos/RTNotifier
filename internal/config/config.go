package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL, RedisAddr, BotToken, ChatID string
	Interval, Cooldown                       time.Duration
}

func Load() Config {
	return Config{DatabaseURL: get("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/payments?sslmode=disable"), RedisAddr: get("REDIS_ADDR", "localhost:6379"), BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"), ChatID: os.Getenv("TELEGRAM_CHAT_ID"), Interval: duration("CHECK_INTERVAL", 15*time.Second), Cooldown: duration("NOTIFY_COOLDOWN", 24*time.Hour)}
}
func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func duration(k string, d time.Duration) time.Duration {
	v, err := time.ParseDuration(os.Getenv(k))
	if err != nil {
		return d
	}
	return v
}
