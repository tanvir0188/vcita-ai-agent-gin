# vcita AI Patient Assistant

A HIPAA-compliant Go service that integrates with the vcita/inTandem platform to provide AI-powered patient communication, medication refill scheduling, and smart appointment booking.

## Architecture

```
vcita Webhooks → HTTPS Server (Gin) → Webhook Handler → AI Service (internal/ai)
                                      → vcita Client (internal/vcita)
                                      → Encrypted SQLite DB (internal/store)
                                      → Audit Logger (internal/audit)
                                      → Notifier (internal/notification)
Cron Scheduler → Refill Reminders (internal/medication) → vcita API + Slack Escalation
```

## Structure Pattern

Each feature module is structured cleanly inside the `internal/` directory:
- `routes.go` — HTTP route registration and handlers.
- `store.go` — Database operations and structures.
- `service.go` — Service logic and third-party integrations (e.g. vcita platform APIs).
- `types.go` — Package-local payload definitions.

Key directories:
- `cmd/api/` — API server setup, routing, and middlewares.
- `cmd/server/` — Main application startup and dependency injection.
- `internal/admin/` — Admin portal sub-modules (`user`, `appointments`, `conversations`, `medicationrefill`).
- `internal/ai/` — OpenAI classification pipelines and FAQ handling.
- `internal/appointment/` — Appointment booking matching and validation logic.
- `internal/medication/` — Refill reminder scheduler and note processing utilities.
- `internal/vcita/` — Shared API client for communication with the vcita platform.
- `internal/notification/` — Slack integration.
- `internal/shared/` — Common HTTP bindings and validators.

## Prerequisites

- Go 1.25+
- A vcita/inTandem business account with API token
- SQLite database DSN
- OpenAI API credentials

## Setup

1. Copy `.env.example` to `.env` and fill in all values.
2. Generate encryption key: `openssl rand -hex 32` (or via `internal/crypto/phi.go` tools).
3. Build the project:
   ```bash
   CGO_ENABLED=1 go build -o vcita-agent ./cmd/server
   ```
4. Run:
   ```bash
   ./vcita-agent
   ```

## HIPAA Compliance Notes

- **AES-256-GCM Encryption**: All PHI data fields are encrypted at-rest.
- **TLS 1.2+ Transit**: All client-server request interfaces are protected.
- **Clean Logs**: Application logs never contain PHI context.
- **Audit Logs**: Events are recorded locally in structured JSON and in the database.
