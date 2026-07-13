package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/admin/middleware"
	"github.com/tanvir0188/vcita-ai-agent/internal/admin/user"
)

func SettingsRoutes(rg *gin.RouterGroup, s *Store) {
	protected := rg.Group("/")
	userStore := user.NewStore(s.db)
	protected.Use(user.WithJWTAuth(userStore))

	protected.GET("/settings", s.GetSettings)
	protected.PATCH("/settings", s.ToggleSettings)
}

type ToggleSettingsPayload struct {
	SystemEnabled bool `json:"system_enabled"`
}

func (s *Store) GetSettings(c *gin.Context) {
	setting, err := s.GetSystemSetting()
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to get system settings",
			Err:     err,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    setting,
	})
}

func (s *Store) ToggleSettings(c *gin.Context) {
	var payload ToggleSettingsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusBadRequest,
			Message: "failed to parse request payload",
			Err:     err,
		})
		c.Abort()
		return
	}

	err := s.UpdateSystemSetting(payload.SystemEnabled)
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update system settings",
			Err:     err,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "system settings updated successfully",
		"data": gin.H{
			"system_enabled": payload.SystemEnabled,
		},
	})
}
