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

	"github.com/tanvir0188/vcita-ai-agent/cmd/api"
	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/crypto"
	"github.com/tanvir0188/vcita-ai-agent/internal/logger"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Envs

	if err := logger.Init(); err != nil {
		return err
	}

	defer logger.Log.Sync()

	enc, err := crypto.NewEncryptor(cfg.PHIEncryptionKey)
	if err != nil {
		return err
	}

	db, err := store.New(cfg.DBDSN, enc, logger.Log)
	if err != nil {
		return err
	}

	defer db.Close()

	auditor, err := audit.New(
		cfg.AuditLogPath,
		db.WriteAuditLog,
	)
	if err != nil {
		return err
	}

	defer auditor.Close()

	server := api.NewAPIServer(
		cfg,
		db,
		auditor,
		logger.Log,
	)

	return server.Run()
}

func (s *APIServer) Run() error {
  router := s.setupRouter()

  srv := &http.Server{
    Addr:         ":" + s.cfg.ServerPort,
    Handler:      router,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  120 * time.Second,
  }

  if s.cfg.TLSCertFile != "" &&
    s.cfg.TLSKeyFile != "" {

    srv.TLSConfig = &tls.Config{
      MinVersion: tls.VersionTLS12,
      CurvePreferences: []tls.CurveID{
        tls.X25519,
        tls.CurveP256,
      },
    }

    go func() {
      if err := srv.ListenAndServeTLS(
        s.cfg.TLSCertFile,
        s.cfg.TLSKeyFile,
      ); err != nil &&
        !errors.Is(err, http.ErrServerClosed) {
        s.log.Fatal("https server error", zap.Error(err))
      }
    }()
  } else {
    go func() {
      if err := srv.ListenAndServe(); err != nil &&
        !errors.Is(err, http.ErrServerClosed) {
        s.log.Fatal("http server error", zap.Error(err))
      }
    }()
  }

  s.auditor.Log(
    "server_started",
    "",
    "system",
    fmt.Sprintf(
      "port=%s tls=%v",
      s.cfg.ServerPort,
      s.cfg.TLSCertFile != "",
    ),
  )

  quit := make(chan os.Signal, 1)

  signal.Notify(
    quit,
    syscall.SIGINT,
    syscall.SIGTERM,
  )

  <-quit

  ctx, cancel := context.WithTimeout(
    context.Background(),
    30*time.Second,
  )
  defer cancel()

  return srv.Shutdown(ctx)
}
