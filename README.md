# RTNotifier

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/) [![CI](https://github.com/igorynos/RTNotifier/actions/workflows/ci.yml/badge.svg)](https://github.com/igorynos/RTNotifier/actions/workflows/ci.yml)

Production-style monitoring worker that analyzes payment-platform data and delivers actionable Telegram alerts. This is the Go successor to [RTNotifier-python](https://github.com/igorynos/RTNotifier-python).

## Alerts

- Provider failure bursts (`process_pending` / `process_failed`)
- Client failure bursts
- Previously active clients whose traffic stopped
- Newly launched clients reaching the transaction threshold
- Merchant/provider pairs with low 24-hour conversion

## Reliability

- PostgreSQL queries run through a bounded connection pool
- Redis `SET NX` provides distributed duplicate suppression and cooldowns
- Trigger, state, and sender interfaces isolate business logic from infrastructure
- Telegram requests carry cancellation and validate HTTP status codes
- Configurable schedules and graceful SIGINT/SIGTERM shutdown

## Layout

```text
cmd/rtnotifier       composition root and scheduler
internal/triggers    payment monitoring queries
internal/notifier    orchestration and interfaces
internal/state       Redis cooldown state
internal/telegram    Telegram Bot API adapter
configs              environment template
deployments          Docker Compose deployment
```

## Run

```bash
cp configs/.env.example .env
docker compose --env-file .env -f deployments/compose.yml up --build
```

The payment database is intentionally external: set `DATABASE_URL` to a schema containing the payment-platform tables documented in the trigger queries.

```bash
make test
make lint
make build
```
