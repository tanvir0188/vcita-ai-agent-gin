package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string

	// Server
	ServerPort  string
	TLSCertFile string
	TLSKeyFile  string

	// vcita
	VcitaAPIBase        string
	VcitaDirectoryToken string
	VcitaBusinessToken  string
	VcitaWebhookSecret  string

	// Database
	DBDriver string
	DBDSN    string

	// Encryption
	PHIEncryptionKey       []byte
	JWTSecret              string
	JWTExpirationInSeconds int

	// Notifications
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	AlertEmailTo string
	OpenAPIKey   string

	// AI Service
	AIServiceURL   string
	AIServiceToken string

	// Audit
	AuditLogPath           string
	SlackMessageWebhookUrl string
}

var Envs = initConfig()

func initConfig() Config {
	_ = godotenv.Load()

	keyBytes := decodeEncryptionKey()

	cfg := Config{
		Environment: getEnv("ENVIRONMENT", "local"),

		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// vcita
		VcitaAPIBase:        mustGetEnv("VCITA_API_BASE"),
		VcitaDirectoryToken: mustGetEnv("VCITA_DIRECTORY_TOKEN"),
		VcitaBusinessToken:  mustGetEnv("VCITA_BUSINESS_TOKEN"),
		VcitaWebhookSecret:  mustGetEnv("VCITA_WEBHOOK_SECRET"),

		// Database
		DBDriver: mustGetEnv("DB_DRIVER"),
		DBDSN:    mustGetEnv("DB_DSN"),

		// Encryption
		PHIEncryptionKey:       keyBytes,
		JWTSecret:              mustGetEnv("JWTSecret"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),

		// Notifications
		SMTPHost:     mustGetEnv("SMTP_HOST"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:     mustGetEnv("SMTP_USER"),
		SMTPPassword: mustGetEnv("SMTP_PASSWORD"),
		SMTPFrom:     mustGetEnv("SMTP_FROM"),
		AlertEmailTo: mustGetEnv("ALERT_EMAIL_TO"),
		OpenAPIKey:   mustGetEnv("OPEN_AI_KEY"),

		// AI Service
		AIServiceURL:   mustGetEnv("AI_SERVICE_URL"),
		AIServiceToken: mustGetEnv("AI_SERVICE_TOKEN"),

		// Audit
		AuditLogPath:           mustGetEnv("AUDIT_LOG_PATH"),
		SlackMessageWebhookUrl: mustGetEnv("SLACK_MESSAGE_WEBHOOK_URL"),
	}

	if cfg.Environment != "local" {
		cfg.TLSCertFile = mustGetEnv("SERVER_TLS_CERT_FILE")
		cfg.TLSKeyFile = mustGetEnv("SERVER_TLS_KEY_FILE")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}

func decodeEncryptionKey() []byte {
	keyHex := mustGetEnv("PHI_ENCRYPTION_KEY")

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		panic("PHI_ENCRYPTION_KEY must be a valid hex string")
	}

	if len(keyBytes) != 32 {
		panic(fmt.Sprintf(
			"PHI_ENCRYPTION_KEY must decode to exactly 32 bytes, got %d",
			len(keyBytes),
		))
	}

	return keyBytes
}
