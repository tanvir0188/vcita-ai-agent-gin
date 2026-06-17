package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
)

func ClientRoutes(rg *gin.RouterGroup, s *Store) {
	rg.GET("/clients", s.ListConversation)

	// rg.PATCH("/clients/:id", s.UpdateAutoMessageSetting)
	// rg.GET("/clients/:id", s.ShowClientDetailPage)
}

func (s *Store) ListConversation(c *gin.Context) {
	var conversations []store.Conversation

	if err := s.db.Find(&conversations).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "base.tmpl", gin.H{
			"title": "Conversations",
			"error": "Failed to retrieve conversations",
		})
		return
	}

	c.HTML(http.StatusOK, "base.tmpl", gin.H{
		"title":         "Conversations",
		"conversations": conversations,
	})
}
