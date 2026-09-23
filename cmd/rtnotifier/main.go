package main

import (
	"context"
	"github.com/igorynos/RTNotifier/internal/notifier"
	"github.com/igorynos/RTNotifier/internal/telegram"
	"log"
	"os"
	"time"
)

type heartbeat struct{}

func (heartbeat) Name() string { return "heartbeat" }
func (heartbeat) Check(context.Context) (bool, string, error) {
	return true, "RTNotifier is running", nil
}
func main() {
	token, chat := os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID")
	if token == "" || chat == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID are required")
	}
	r := notifier.NewRunner(telegram.Client{Token: token, ChatID: chat}, 24*time.Hour, heartbeat{})
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	ctx := context.Background()
	for {
		if err := r.RunOnce(ctx); err != nil {
			log.Printf("notification cycle: %v", err)
		}
		<-ticker.C
	}
}
