package medicationrefill

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/user"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
)

func MedicationRemindertRoutes(rg *gin.RouterGroup, s *Store) {
	protected := rg.Group("/")
	reminderStore := user.NewStore(s.db)
	protected.Use(user.WithJWTAuth(reminderStore))
	protected.GET("/medication-reminders", s.ListMedicationReminders)
	protected.PATCH("/medication-reminders/:id", s.ToggleMedicationReminder)

}

func (s *Store) ListMedicationReminders(c *gin.Context) {
	var medication []store.ClientSyncState

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

	if err := s.db.Model(&store.ClientSyncState{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve medication reminders",
		})
		return
	}
	if err := s.db.Limit(limit).Offset(offset).Find(&medication).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve medication reminders",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    medication,
		"meta": gin.H{
			"limit": limit,
			"page":  page,
			"pages": total / int64(limit),
			"total": total,
		},
	})
}

func (s *Store) ToggleMedicationReminder(c *gin.Context) {
	id := c.Param("id")
	err := s.ToggleReminder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to toggle medication reminder",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "medication reminder toggled successfully",
	})
}
