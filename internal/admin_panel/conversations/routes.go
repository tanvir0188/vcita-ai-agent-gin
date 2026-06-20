package client

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/user"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
)

func ClientRoutes(rg *gin.RouterGroup, s *Store) {
	protected := rg.Group("/")
	userStore := user.NewStore(s.db)
	protected.Use(user.WithJWTAuth(userStore))
	protected.GET("/conversations", s.ListConversation)

	// rg.PATCH("/clients/:id", s.UpdateAutoMessageSetting)
	// rg.GET("/clients/:id", s.ShowClientDetailPage)
}

func (s *Store) ListConversation(c *gin.Context) {
	var conversations []store.Conversation

	log.Println("ListConversation called")

	if err := s.db.Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve conversations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": conversations,
	})
}
