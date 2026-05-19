// Command server is the entry point for the vcita AI patient assistant.
// Gin replaces chi for HTTP routing; GORM replaces sqlx for database access.
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/crypto"
	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"github.com/tanvir0188/vcita-ai-agent/internal/notify"
	"github.com/tanvir0188/vcita-ai-agent/internal/scheduler"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
	"github.com/tanvir0188/vcita-ai-agent/internal/webhook"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// ── 1. Config ─────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	_ = cfg

	err = logger.Init()
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}

	defer logger.Log.Sync() //nolint:errcheck

	logger.Log.Info("vcita-ai-agent starting")

	// ── 3. Encryption (AES-256-GCM) ───────────────────────────────────────────
	enc, err := crypto.NewEncryptor(cfg.PHIEncryptionKey)
	if err != nil {
		return fmt.Errorf("crypto: %w", err)
	}

	// ── 4. Database (GORM + SQLite) ───────────────────────────────────────────
	if err := os.MkdirAll("data", 0700); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}
	db, err := store.New(cfg.DBDSN, enc, logger.Log)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer db.Close()
	logger.Log.Info("database connected and migrated")

	// ── 5. Audit logger ───────────────────────────────────────────────────────
	if err := os.MkdirAll("logs", 0700); err != nil {
		return fmt.Errorf("create logs dir: %w", err)
	}
	auditor, err := audit.New(cfg.AuditLogPath, db.WriteAuditLog)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}
	defer auditor.Close()

	// ── 6. vcita API client ───────────────────────────────────────────────────
	vcitaClient := vcita.NewAPIClient(cfg.VcitaAPIBase, cfg.VcitaBusinessToken)

	// ── 7. Notifier ───────────────────────────────────────────────────────────
	notifier := notify.NewNotifier(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom, cfg.AlertEmailTo, logger.Log)

	// ── 9. Webhook handler ────────────────────────────────────────────────────
	wh := webhook.New(cfg.VcitaWebhookSecret, db, vcitaClient, auditor, logger.Log)
	conversationReadHandler := webhook.NewConversation(cfg.VcitaWebhookSecret, db, auditor, logger.Log)

	// ── 10. Gin router ────────────────────────────────────────────────────────
	// Set Gin to release mode in production — disables debug noise in logs
	gin.SetMode(gin.ReleaseMode)

	r := gin.New() // gin.New() instead of gin.Default() so we control all middleware

	// Middleware stack
	r.Use(ginZapLogger(logger.Log)) // structured access log, no body logging
	r.Use(gin.Recovery())           // recover from panics
	r.Use(securityHeaders())        // HSTS, X-Frame-Options, etc.

	// Routes
	r.POST("/webhook", wh.Handle)
	r.POST("/webhook/conversation-read", conversationReadHandler.ConversationReadHandle)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ── 11. Scheduler ─────────────────────────────────────────────────────────
	sched := scheduler.New(db, vcitaClient, notifier, auditor, logger.Log)
	if err := sched.Start(); err != nil {
		return fmt.Errorf("scheduler: %w", err)
	}
	defer sched.Stop()

	// ── 12. HTTP(S) server ────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		srv.TLSConfig = &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		}
		logger.Log.Info("starting HTTPS server", zap.String("addr", srv.Addr))
		go func() {
			if err := srv.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.Fatal("HTTPS server error", zap.Error(err))
			}
		}()
	} else {
		logger.Log.Warn("TLS not configured – HTTP only (development mode)")
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.Fatal("HTTP server error", zap.Error(err))
			}
		}()
	}

	auditor.Log("server_started", "", "system",
		fmt.Sprintf("port=%s tls=%v", cfg.ServerPort, cfg.TLSCertFile != ""))

	// ── 13. Graceful shutdown ─────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down…")
	auditor.Log("server_shutdown", "", "system", "signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	logger.Log.Info("shutdown complete")
	return nil
}

// ── Middleware ────────────────────────────────────────────────────────────────

// ginZapLogger creates a Gin middleware that logs each request with zap.
// It intentionally omits request bodies and query strings to avoid PHI leakage.
func ginZapLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("http",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path), // path only — no query params
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// securityHeaders adds HIPAA/security-relevant HTTP response headers.
func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Next()
	}
}

// ── Logger ────────────────────────────────────────────────────────────────────

func buildLogger() (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = zapcore.DebugLevel
	}
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	return cfg.Build()
}
