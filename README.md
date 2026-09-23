# RTNotifier

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/) [![CI](https://github.com/igorynos/RTNotifier/actions/workflows/ci.yml/badge.svg)](https://github.com/igorynos/RTNotifier/actions/workflows/ci.yml)

Production-style monitoring worker that analyzes payment-platform data and delivers actionable Telegram alerts. This is the Go successor to [RTNotifier-python](https://github.com/igorynos/RTNotifier-python).

## 🚨 Alert Scenarios

- 🏦 **Provider incidents:** Detects bursts of `process_pending` and `process_failed` transactions.
- 👤 **Client incidents:** Finds commerce accounts accumulating unsuccessful operations.
- ⛔ **Traffic stops:** Reports clients that were active during the last 24 hours but became silent.
- 🚀 **Client launches:** Identifies new clients reaching the configured successful-transaction threshold.
- 📉 **Low conversion:** Calculates 24-hour conversion by merchant and provider and reports degraded routes.

## 🛡️ Reliability

- 🗄️ PostgreSQL queries execute through a bounded connection pool.
- 🔒 Redis `SET NX` provides atomic distributed duplicate suppression.
- ⏱️ Notification intervals and cooldown windows are configured independently.
- 🧩 Trigger, state, and sender interfaces separate business rules from infrastructure.
- 🌐 Telegram requests carry context cancellation and validate response status codes.
- 📴 SIGINT and SIGTERM stop the worker without abandoning the active process abruptly.

## 🔄 Processing Flow

```text
Scheduler
   │
   ▼
PostgreSQL trigger query ──► incident detected?
                                  │
                                  ▼
                         Redis cooldown gate
                                  │
                                  ▼
                         Telegram Bot API
```

## 🏗️ Project Layout

```text
cmd/rtnotifier       composition root and scheduler
internal/triggers    payment monitoring queries
internal/notifier    orchestration and interfaces
internal/state       Redis cooldown state
internal/telegram    Telegram Bot API adapter
configs              environment template
deployments          Docker Compose deployment
```

## 🚀 Run

```bash
cp configs/.env.example .env
docker compose --env-file .env -f deployments/compose.yml up --build
```

The payment database is intentionally external: set `DATABASE_URL` to a schema containing the payment-platform tables documented in the trigger queries.

## ⚙️ Configuration

| Variable | Purpose | Default |
| :--- | :--- | :--- |
| `DATABASE_URL` | Payment PostgreSQL connection | Local PostgreSQL |
| `REDIS_ADDR` | Distributed notification state | `localhost:6379` |
| `TELEGRAM_BOT_TOKEN` | Telegram bot credentials | Required |
| `TELEGRAM_CHAT_ID` | Destination channel or group | Required |
| `CHECK_INTERVAL` | Trigger evaluation frequency | `15s` |
| `NOTIFY_COOLDOWN` | Duplicate notification window | `24h` |

## 🧪 Development

```bash
make test
make lint
make build
```

## 🐍 Previous Implementation

The original Python service remains available as [RTNotifier-python](https://github.com/igorynos/RTNotifier-python).
