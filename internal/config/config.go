// Package config loads and validates all runtime configuration from environment variables.
// No sensitive values are ever hardcoded. All PHI-adjacent secrets are required at startup.
package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the complete application configuration.
// All fields are validated during Load(); missing required values cause a fatal error.
type Config struct {
	Environment string
	// Server
	ServerPort  string
	TLSCertFile string
	TLSKeyFile  string

	// vcita / inTandem
	VcitaAPIBase        string
	VcitaDirectoryToken string
	VcitaBusinessToken  string
	VcitaWebhookSecret  string

	// Database
	DBDriver string
	DBDSN    string

	// Encryption – 32-byte AES-256 key decoded from hex env var
	PHIEncryptionKey []byte

	// Notifications
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	AlertEmailTo string

	// AI Service (internal microservice stub)
	AIServiceURL   string
	AIServiceToken string

	// Audit
	AuditLogPath string
}

// Load reads configuration from environment variables (and optionally a .env file).
// Returns an error if any required variable is absent or invalid.
func Load() (*Config, error) {
	// Load .env if present (ignored in production where env vars are injected by systemd/Docker)
	_ = godotenv.Load()
	var Environment string = os.Getenv("ENVIRONMENT")
	fmt.Print(Environment)

	cfg := &Config{}

	// ── Required fields ──────────────────────────────────────────────────────
	required := map[string]*string{
		"ENVIRONMENT":           &cfg.Environment,
		"SERVER_PORT":           &cfg.ServerPort,
		"VCITA_API_BASE":        &cfg.VcitaAPIBase,
		"VCITA_DIRECTORY_TOKEN": &cfg.VcitaDirectoryToken,
		"VCITA_BUSINESS_TOKEN":  &cfg.VcitaBusinessToken,
		"VCITA_WEBHOOK_SECRET":  &cfg.VcitaWebhookSecret,
		"DB_DRIVER":             &cfg.DBDriver,
		"DB_DSN":                &cfg.DBDSN,
		"SMTP_HOST":             &cfg.SMTPHost,
		"SMTP_USER":             &cfg.SMTPUser,
		"SMTP_PASSWORD":         &cfg.SMTPPassword,
		"SMTP_FROM":             &cfg.SMTPFrom,
		"ALERT_EMAIL_TO":        &cfg.AlertEmailTo,
		"AI_SERVICE_URL":        &cfg.AIServiceURL,
		"AI_SERVICE_TOKEN":      &cfg.AIServiceToken,
		"AUDIT_LOG_PATH":        &cfg.AuditLogPath,
	}
	for key, dest := range required {
		val := os.Getenv(key)
		if val == "" {
			return nil, fmt.Errorf("required environment variable %q is not set", key)
		}
		*dest = val
	}

	// TLS cert/key – required only in production (skip if SERVER_PORT == 8080 for local dev)
	cfg.Environment = os.Getenv("ENVIRONMENT")

	cfg.TLSCertFile = os.Getenv("SERVER_TLS_CERT_FILE")
	cfg.TLSKeyFile = os.Getenv("SERVER_TLS_KEY_FILE")
	if cfg.Environment == "local" {
		cfg.TLSCertFile = ""
		cfg.TLSKeyFile = ""
	}

	// ── SMTP port ────────────────────────────────────────────────────────────
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("SMTP_PORT must be an integer, got %q", portStr)
	}
	cfg.SMTPPort = port

	// ── PHI encryption key ───────────────────────────────────────────────────
	keyHex := os.Getenv("PHI_ENCRYPTION_KEY")
	if keyHex == "" {
		return nil, fmt.Errorf("required environment variable PHI_ENCRYPTION_KEY is not set")
	}
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("PHI_ENCRYPTION_KEY must be a valid hex string: %w", err)
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("PHI_ENCRYPTION_KEY must decode to exactly 32 bytes (AES-256), got %d", len(keyBytes))
	}
	cfg.PHIEncryptionKey = keyBytes
	fmt.Println("Configuration loaded successfully")

	return cfg, nil
}
