# vcita-ai-agent-gin Project Overview and Development Sequence

## Project Description

This is a HIPAA-compliant Go service that integrates with the vcita/inTandem platform to provide AI-powered patient communication, medication refill scheduling, and smart appointment booking. It uses Gin for HTTP routing, GORM for database access, and AES-256-GCM for PHI encryption.

## Architecture Overview

```
vcita Webhooks → HTTPS Server (Gin) → Webhook Handler → AI Service Interface
                                      → vcita API Client
                                      → Encrypted SQLite DB (GORM)
                                      → Audit Logger (File + DB)
                                      → Notifier (SMTP)
Cron Scheduler → Refill Reminders → vcita API + Staff Email
```

## File Creation and Dependency Sequence

The project was developed following a logical sequence to ensure proper layering and dependencies. Here's the order in which components were created and how they depend on each other:

### 1. Core Infrastructure (`internal/` packages)

#### `internal/config/config.go`
- **Purpose**: Centralized configuration loading from environment variables
- **Dependencies**: None (pure stdlib)
- **Why first**: All other components need configuration to initialize
- **Key functions**:
  - `Load()`: Reads env vars, validates required fields, decodes hex encryption key

#### `internal/crypto/phi.go`
- **Purpose**: HIPAA-compliant AES-256-GCM encryption for PHI data
- **Dependencies**: `config` (for encryption key)
- **Why next**: Database and other components need encryption before storing PHI
- **Key functions**:
  - `NewEncryptor(key []byte)`: Creates encryptor instance
  - `Encrypt()/Decrypt()`: Symmetric encryption methods
  - `GenerateKey()`: Generates new random 32-byte key

#### `internal/store/`
- **Purpose**: Database layer with GORM and encrypted PHI storage
- **Dependencies**: `crypto` (for PHI encryption), `config` (for DSN)
- **Why next**: Core data persistence needed before business logic
- **Key files**:
  - `db.go`: Database connection and migrations
  - `models.go`: GORM model definitions
  - `conversations.go`: Message/conversation storage
  - `medications.go`: Medication and refill data
  - `audit.go`: Audit log storage

#### `internal/audit/logger.go`
- **Purpose**: HIPAA audit logging (file + database)
- **Dependencies**: `store` (for DB audit writes)
- **Why next**: Audit logging needed for compliance in all operations
- **Key functions**:
  - `New()`: Creates logger with file and optional DB sink
  - `Log()`: Records audit events without PHI

#### `internal/vcita/`
- **Purpose**: HTTP client for vcita/inTandem API
- **Dependencies**: None (pure HTTP client)
- **Why next**: External API integration needed for webhook responses
- **Key files**:
  - `types.go`: API request/response structs
  - `client.go`: HTTP client methods (SendMessage, GetClient, etc.)

#### `internal/notify/notifier.go`
- **Purpose**: SMTP-based staff alerts for escalations and reminders
- **Dependencies**: None (pure SMTP client)
- **Why next**: Notification system needed for automated alerts
- **Key functions**:
  - `NewNotifier()`: SMTP configuration
  - `AlertStaff()`: Sends structured alerts
  - `SendRefillReminder()`: Medication refill notifications

#### `internal/ai/stub.go`
- **Purpose**: AI service interface placeholder
- **Dependencies**: None (interface definition)
- **Why next**: Interface needed before webhook handler
- **Key interface**:
  - `AIService`: ProcessMessage, ExtractMedication, SuggestAppointment

#### `internal/webhook/handler.go`
- **Purpose**: HTTP webhook endpoint processing
- **Dependencies**: `vcita`, `store`, `ai`, `audit`, `notify`
- **Why next**: Main business logic orchestrator
- **Key functions**:
  - `New()`: Handler constructor
  - `Handle()`: Main webhook processing

#### `internal/scheduler/refill.go`
- **Purpose**: Cron-based medication refill reminders
- **Dependencies**: `store`, `vcita`, `notify`, `audit`
- **Why next**: Background job processing
- **Key functions**:
  - `New()`: Scheduler with cron jobs
  - `processOne()`: Individual refill processing

### 2. Entry Point (`cmd/server/main.go`)
- **Purpose**: Application bootstrap and orchestration
- **Dependencies**: All `internal/` packages
- **Why last**: Ties everything together
- **Key sequence**:
  1. Load config
  2. Initialize logger
  3. Create encryptor
  4. Connect database
  5. Setup audit logger
  6. Create vcita client
  7. Setup notifier
  8. Initialize AI service
  9. Create webhook handler
  10. Setup Gin router
  11. Start scheduler
  12. Start HTTPS server
  13. Graceful shutdown handling

### 3. Scripts (`scripts/`)
- **Purpose**: Utility scripts for deployment and setup
- **Dependencies**: Various `internal/` packages
- **Why parallel**: Utilities developed as needed
- **Key scripts**:
  - `genkey/main.go`: Generate encryption keys
  - `register-webhooks/main.go`: Setup vcita webhooks (stub)
  - `deploy.sh`: Deployment script
  - `vcita-agent.service`: Systemd service file

### 4. Infrastructure (`Dockerfile`, `go.mod`)
- **Purpose**: Containerization and dependency management
- **Dependencies**: All source code
- **Why throughout**: Updated as components added
- **Key aspects**:
  - Multi-stage Docker build
  - Go module dependencies (Gin, GORM, Zap, etc.)
  - CGO enabled for SQLite

## Development Flow

1. **Configuration First**: Environment-based config ensures no hardcoded secrets
2. **Security Layer**: Encryption initialized early for PHI protection
3. **Data Layer**: Database setup before business logic
4. **External APIs**: vcita client for platform integration
5. **Communication**: Notification system for staff alerts
6. **AI Integration**: Pluggable interface for ML services
7. **HTTP Layer**: Webhook processing as main entry point
8. **Background Jobs**: Scheduler for automated tasks
9. **Orchestration**: Main function coordinates all components
10. **Deployment**: Scripts and containers for production

## Key Design Decisions

- **Dependency Injection**: All components receive dependencies explicitly
- **Interface Segregation**: AI service as interface allows easy swapping
- **Layered Architecture**: Clear separation between HTTP, business logic, data
- **Security First**: PHI encryption, audit logging, TLS everywhere
- **Observability**: Structured logging, health checks, graceful shutdown
- **Testability**: Pure functions, interface mocks, minimal side effects

This sequence ensures each component has its dependencies ready and follows Go best practices for maintainable, testable code.