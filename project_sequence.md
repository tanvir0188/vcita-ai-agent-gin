# vcita-ai-agent-gin Project Overview and Development Sequence

## Project Description

This is a HIPAA-compliant Go service that integrates with the vcita/inTandem platform to provide AI-powered patient communication, medication refill scheduling, and smart appointment booking. It uses Gin for HTTP routing, GORM for database access, and AES-256-GCM for PHI encryption.

## Architecture Overview

```
vcita Webhooks → HTTPS Server (Gin) → Webhook Handler → AI Service (internal/ai)
                                      → vcita Client (internal/vcita)
                                      → Encrypted SQLite DB (internal/store)
                                      → Audit Logger (internal/audit)
                                      → Notifier (internal/notification)
Cron Scheduler → Refill Reminders (internal/medication) → vcita API + Slack Escalation
```

## File Creation and Dependency Sequence

The project is organized in a modular structure where each feature package follows the routes-store-service-types architecture.

### 1. Core Infrastructure (`internal/` packages)

#### `internal/config/config.go`
- **Purpose**: Centralized configuration loading from environment variables.
- **Dependencies**: None.
- **Key functions**:
  - `initConfig()`: Reads env vars, validates required fields, decodes hex encryption key.

#### `internal/crypto/phi.go`
- **Purpose**: HIPAA-compliant AES-256-GCM encryption for PHI data.
- **Dependencies**: `config` (for encryption key).
- **Key functions**:
  - `NewEncryptor(key []byte)`: Creates encryptor instance.
  - `Encrypt()/Decrypt()`: Symmetric encryption/decryption helper methods.

#### `internal/store/`
- **Purpose**: Central database layer with GORM model mappings and encryption.
- **Dependencies**: `crypto`, `config`.
- **Key files**:
  - `db.go`: Database connection pool configuration and auto-migrations.
  - `models.go`: Core GORM structural models (`User`, `Conversation`, `ClientSyncState`, `Appointment`, `AuditLog`).
  - `service.go`: CRUD operations for conversation message creation.
  - `audit.go`: Appends HIPAA-compliant audit records.

#### `internal/audit/logger.go`
- **Purpose**: Append-only HIPAA audit logging (file + database).
- **Dependencies**: `store` (for DB audit writes).
- **Key functions**:
  - `New()`: Creates logger with file and optional DB sink.
  - `Log()`: Records audit events without PHI details.

#### `internal/vcita/`
- **Purpose**: Central shared client for vcita/inTandem platform APIs.
- **Dependencies**: `config`, `logger`.
- **Key files**:
  - `client.go`: Generic HTTP helper functions (`DoRequest`, `DoRequestWithBody`).
  - `message.go`: API actions for message routing (`CreateMessage`, `GetMessageHistory`, `GetSmartReplyEmail`).
  - `contact.go`: API actions for contacts (`GetClientDetail`, `GetClientNotes`).
  - `staff.go`: API actions for staffs (`GetStaffList`, `GetStaffIDs`).
  - `types.go`: Shared response structures.

#### `internal/shared/`
- **Purpose**: Cross-cutting HTTP helper functions used by Gin routes.
- **Dependencies**: `store`.
- **Key functions**:
  - `BindAndValidate()`: Validates and binds incoming request payloads.
  - `GetTokenFromRequest()`: Extracts access tokens from authorization headers.
  - `GetUserFromRequest()`: Retrieves the current authenticated user context.

#### `internal/notification/`
- **Purpose**: Slack Block Kit alerts for human escalations.
- **Dependencies**: `vcita`, `logger`.
- **Key functions**:
  - `SendMessageToSlack()`: Formats and sends notifications to Slack.

---

### 2. Feature Modules (`internal/` packages)

#### `internal/ai/`
- **Purpose**: OpenAI AI evaluation and classification logic.
- **Dependencies**: None.
- **Key files**:
  - `open_ai.go`: Handles core OpenAI classification (`HumanInterventionNeeded`, `ConfirmedAppointmentDate`, `PredictMedicationRefill`).
  - `systemPrompt.go`: Prompts utilized during execution.
  - `readDoc.go`: FAQ base extracted from medical documents.
  - `types.go`: Local AI payload mappings.

#### `internal/appointment/`
- **Purpose**: Appointment scheduling and slots orchestration.
- **Dependencies**: `ai`, `store`, `vcita`, `config`.
- **Key files**:
  - `handler.go`: Handles appointment slot verification and schedules bookings.
  - `service.go`: Interacts with vcita platform API availability endpoints.
  - `types.go`: Slot booking payloads.

#### `internal/medication/`
- **Purpose**: Background tasks for medication refill reminders.
- **Dependencies**: `ai`, `store`, `vcita`, `config`.
- **Key files**:
  - `scheduler.go`: Manages cron intervals and triggers reminder broadcasts.
  - `service.go`: Queries patient lists and checks current notes on vcita.
  - `utils.go`: Sanitization helper utilities (HTML to markdown, notes filters).

#### `internal/webhook/`
- **Purpose**: Webhook listener endpoints for message triggers.
- **Dependencies**: `ai`, `store`, `vcita`, `appointment`, `notification`.
- **Key files**:
  - `utils.go`: Processes client webhook triggers, handles wait periods, and schedules smart replies or escalations.
  - `types.go`: Webhook structures.

#### `internal/admin/`
- **Purpose**: Backoffice administrator portal REST endpoints.
- **Dependencies**: `store`, `shared`, `config`.
- **Key packages**:
  - `user/`: Handles login, registration (`service.go` looks up staff by email from `vcita`), and password modifications.
  - `appointments/`: Lists paginated scheduled appointments.
  - `conversations/`: Lists user conversations and controls auto-replies.
  - `medicationrefill/`: Manages automated daily reminders.

---

### 3. Application Entry Points

#### `cmd/api/`
- **Purpose**: APIServer constructor and routing table definition.
- **Dependencies**: All feature packages (`admin`, `webhook`, `medication`, `store`, `audit`).
- **Key files**:
  - `api.go`: Gin router initialization, global middleware stack, and endpoint group declarations.

#### `cmd/server/main.go`
- **Purpose**: Service entrypoint.
- **Dependencies**: `cmd/api`, `logger`, `crypto`, `store`, `config`.
- **Key sequence**:
  1. Load environment parameters.
  2. Instantiate core logger.
  3. Instantiate AES encryptor.
  4. Establish store connection.
  5. Setup audit tracker.
  6. Bootstrap daily scheduler.
  7. Start Gin API server.

---

## Key Design Decisions

- **routes-store-service-types Pattern**: Consistently modularizes logical groups to ease maintenance.
- **Decoupled API Actions**: Relocated third-party integrations (e.g., `vcita`) out of general helpers and into isolated modules.
- **Strict Separation of Concerns**: Shared helpers reside in `internal/shared`, and external platform calls are grouped inside `internal/vcita`.
- **Security-First**: Enforces symmetric database encryption, clean audit logging, and parameterized operations.