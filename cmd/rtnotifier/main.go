package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/igorynos/RTNotifier/internal/config"
	"github.com/igorynos/RTNotifier/internal/notifier"
	"github.com/igorynos/RTNotifier/internal/state"
	"github.com/igorynos/RTNotifier/internal/telegram"
	"github.com/igorynos/RTNotifier/internal/triggers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	if cfg.BotToken == "" || cfg.ChatID == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are required")
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cache := state.New(cfg.RedisAddr)
	defer cache.Close()
	all := triggers.All(db)
	items := make([]notifier.Trigger, len(all))
	for i := range all {
		items[i] = all[i]
	}
	runner := notifier.NewRunner(telegram.Client{Token: cfg.BotToken, ChatID: cfg.ChatID}, cfg.Cooldown, items...)
	runner.State = cache
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for {
		if err = runner.RunOnce(ctx); err != nil {
			log.Printf("notification cycle: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
