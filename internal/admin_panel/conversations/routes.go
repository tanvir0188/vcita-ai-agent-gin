package conversations

import (
	"fmt"
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
	protected.PATCH("/conversations/:id", s.toggleAutoReply)

	// rg.PATCH("/clients/:id", s.UpdateAutoMessageSetting)
	// rg.GET("/clients/:id", s.ShowClientDetailPage)
}

func (s *Store) toggleAutoReply(c *gin.Context) {
	id := c.Param("id")

	var convo store.Conversation

	err := s.db.Where("id = ?", id).First(&convo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "conversation not found",
		})
		return
	}

	// toggle boolean field (adjust field name if different)
	convo.HasEscalated = !convo.HasEscalated

	err = s.db.Save(&convo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update conversation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Escalation handled",
		"data": gin.H{
			"id":                 convo.ID,
			"conversation_id":    convo.ConversationID,
			"escalation_handled": convo.HasEscalated,
		},
	})
}

func (s *Store) ListConversation(c *gin.Context) {
	var conversations []store.Conversation

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64

	if err := s.db.Model(&store.Conversation{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count conversations",
		})
		return
	}

	if err := s.db.
		Limit(limit).
		Offset(offset).
		Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve conversations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": conversations,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}
