# RTNotifier

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/) [![CI](https://github.com/igorynos/RTNotifier/actions/workflows/ci.yml/badge.svg)](https://github.com/igorynos/RTNotifier/actions/workflows/ci.yml)

Extensible monitoring service that evaluates business triggers and delivers real-time alerts through Telegram. Rules are independent from delivery infrastructure, so database checks can be tested without calling external APIs.

This is the Go successor to [RTNotifier-python](https://github.com/igorynos/RTNotifier-python).

## Highlights

- Small `Trigger` and `Sender` interfaces
- Periodic background evaluation
- Per-trigger cooldown and duplicate suppression
- Telegram Bot API adapter with context cancellation
- Error propagation with trigger-level context
- Deterministic unit tests for delivery behavior
- Minimal multi-stage Docker image

## Processing flow

```text
Data source → Trigger evaluation → Cooldown check → Telegram sender
```

## Run

```bash
go test ./...
TELEGRAM_BOT_TOKEN=... \
TELEGRAM_CHAT_ID=... \
go run ./cmd/rtnotifier
```

## Roadmap

- PostgreSQL-backed business triggers
- Redis persistence for cooldown state
- Worker pool and graceful shutdown
- Prometheus metrics and structured logging
- Multiple notification destinations
