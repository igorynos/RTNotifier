# RTNotifier

Go rewrite of [RTNotifier-python](https://github.com/igorynos/RTNotifier-python), a configurable monitoring and Telegram notification service.

The current migration provides a trigger engine, cooldown-based duplicate protection, a Telegram Bot API adapter, tests, and a minimal Docker image. Database-backed triggers can implement the small `Trigger` interface without coupling business rules to Telegram.

```bash
go test ./...
TELEGRAM_BOT_TOKEN=... TELEGRAM_CHAT_ID=... go run ./cmd/rtnotifier
```

Next: PostgreSQL trigger repository, Redis-backed cooldown state, Prometheus metrics, and graceful worker orchestration.
