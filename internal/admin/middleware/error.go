package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"error,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// ErrorHandler catches errors added to the context and returns them uniformly.
func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		if c.Writer.Written() {
			return
		}

		err := c.Errors.Last()

		var appErr *AppError
		if errors.As(err.Err, &appErr) {
			log.Error(
				"request failed",
				zap.Error(appErr.Err),
				zap.Int("status", appErr.Status),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.JSON(appErr.Status, gin.H{
				"message": appErr.Message,
			})
			return
		}

		log.Error(
			"request failed",
			zap.Error(err.Err),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
	}
}
