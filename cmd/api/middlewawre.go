package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/webhook"
)

func (s *APIServer) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(ZapLogger(s.log))
	r.Use(gin.Recovery())
	r.Use(securityHeaders())

	api := r.Group("/api/v1")

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
	api.POST("/webhook/conversation-read",
		conversationHandler.ConversationReadHandle)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return r
}
