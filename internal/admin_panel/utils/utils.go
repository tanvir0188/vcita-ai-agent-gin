package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
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

func BindAndValidate(c *gin.Context, payload interface{}) error {
	if err := c.ShouldBindJSON(payload); err != nil {
		return err
	}

	if err := Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return fmt.Errorf("invalid payload: %v", validationErrors)
	}

	return nil
}

func GetUserFromRequest(c *gin.Context) (*store.User, error) {
	userValue, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("unauthorized")
	}

	user, ok := userValue.(store.User)
	if !ok {
		return nil, fmt.Errorf("invalid user context")
	}

	return &user, nil
}
