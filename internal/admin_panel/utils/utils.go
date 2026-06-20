package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func GetTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")

	if authHeader != "" {
		const prefix = "Bearer "

		if len(authHeader) > len(prefix) &&
			authHeader[:len(prefix)] == prefix {
			return authHeader[len(prefix):]
		}

		return authHeader
	}

	return c.Query("token")
}
