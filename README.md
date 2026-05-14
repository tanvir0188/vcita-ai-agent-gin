# vcita AI Patient Assistant

A HIPAA-compliant Go service that integrates with the vcita/inTandem platform to
provide AI-powered patient communication, medication refill scheduling, and smart
appointment booking.

## Architecture

```
vcita Webhooks → HTTPS Server → Webhook Handler → AI Service (interface)
                                               → vcita API Client
                                               → Encrypted SQLite/Postgres
                                               → Audit Logger
Cron Scheduler (daily 08:00)   → Refill Reminders → vcita + Staff Email
```

## Prerequisites

- Go 1.22+
- A vcita/inTandem business account with API token
- VPS with a domain and TLS certificate (Let's Encrypt)
- SMTP credentials (SendGrid recommended)

## Setup

1. Copy `.env.example` to `.env` and fill in all values
2. Generate encryption key: `openssl rand -hex 32`
3. Build: `CGO_ENABLED=1 go build -o vcita-agent ./cmd/server`
4. Run: `./vcita-agent`

## HIPAA Compliance Notes

- All PHI (message content, notes, medication data) is encrypted at rest using AES-256-GCM
- All data in transit uses TLS 1.2+ minimum
- Audit logs record every system action without storing PHI content
- No PHI appears in application logs – only event types, IDs, and actions
- Database files are stored with 0700 permissions
- Service runs as a non-root user

## AI Service Integration

The `internal/ai/stub.go` file contains a placeholder `Stub` that satisfies the
`webhook.AIService` interface. The AI developer should implement this interface
in a separate package and swap it in `cmd/server/main.go`.

## Webhook Events Handled

| Event | Action |
|---|---|
| `message.client_sent_message` | AI response pipeline |
| `message.business_sent_message` | Audit log only |
| `client.updated` | Medication extraction if note type is `current_medications` |
| `appointment.requested` | AI slot suggestion + booking |
