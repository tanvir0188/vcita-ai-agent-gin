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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	client "github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/conversations"
	"github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/user"
	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/webhook"
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

func (s *APIServer) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(ZapLogger(s.log))
	r.Use(gin.Recovery())
	r.Use(securityHeaders())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	//get run time path for loading templates
	// execPath, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }
	// templatePath := filepath.Join(execPath, "templates", "*.tmpl")
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.tmpl")
	r.LoadHTMLGlob("templates/*/*.tmpl")

	public := r.Group("/")

	public.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	store := user.NewStore(s.db.GetGorm())

	clientStore := client.NewStore(s.db.GetGorm())

	api := r.Group("/api/v1")

	user.RegisterRoutes(api, store)
	client.ClientRoutes(api, clientStore)

	webhookHandler := webhook.New(
		s.cfg.VcitaWebhookSecret,
		s.db,
		s.cfg.SlackMessageWebhookUrl,
		s.cfg.OpenAPIKey,
		s.auditor,
		s.log,
	)

	conversationHandler := webhook.NewConversation(
		s.cfg.VcitaWebhookSecret,
		s.db,
		s.auditor,
		s.log,
	)

	api.POST("/webhook", webhookHandler.ConversationCreateHandle)
	r.POST("/webhook/conversation-read",
		conversationHandler.ConversationReadHandle)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return r
}

func (s *APIServer) Run() error {
	router := s.setupRouter()
	s.log.Info("====================================")
	s.log.Info("vcita-ai-agent started successfully")
	s.log.Info(
		"server config",
		zap.String("port", s.cfg.ServerPort),
		zap.Bool("tls", s.cfg.TLSCertFile != ""),
		zap.String("environment", s.cfg.Environment),
	)

	for _, route := range router.Routes() {
		s.log.Info(
			"route",
			zap.String("method", route.Method),
			zap.String("path", route.Path),
		)
	}

	s.log.Info("====================================")

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
