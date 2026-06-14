package api

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

	"go.uber.org/zap"

	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
)

type APIServer struct {
	cfg     config.Config
	db      *store.DB
	auditor *audit.Logger
	log     *zap.Logger
}

func NewAPIServer(
	cfg config.Config,
	db *store.DB,
	auditor *audit.Logger,
	log *zap.Logger,
) *APIServer {
	return &APIServer{
		cfg:     cfg,
		db:      db,
		auditor: auditor,
		log:     log,
	}
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
