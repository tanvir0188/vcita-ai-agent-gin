package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/shared"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
)

const UserKey = "user_id"

func WithJWTAuth(s *Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := shared.GetTokenFromRequest(c)

		if tokenString == "" {
			permissionDenied(c)
			return
		}

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(c)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(c)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("failed to parse claims")
			permissionDenied(c)
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			log.Println("missing sub claim")
			permissionDenied(c)
			return
		}

		userID, err := strconv.Atoi(sub)
		if err != nil {
			log.Printf("failed to convert user id: %v", err)
			permissionDenied(c)
			return
		}

		var u store.User

		err = s.DB().
			Where("id = ?", userID).
			First(&u).
			Error

		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(c)
			return
		}

		c.Set(UserKey, u.ID)
		c.Set("user", u)

		c.Next()
	}
}

func CreateRefreshToken(secret []byte, userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  strconv.Itoa(int(userID)),
		"type": "refresh",
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})

	return token.SignedString(secret)
}

func CreateAccessToken(secret []byte, userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  strconv.Itoa(int(userID)),
		"type": "access",
		"exp":  time.Now().Add(15 * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})

	return token.SignedString(secret)
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(
				"unexpected signing method: %v",
				token.Header["alg"],
			)
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error": "permission denied",
	})
}

func GetUserIDFromContext(c *gin.Context) int {
	userID, exists := c.Get(UserKey)
	if !exists {
		return -1
	}

	switch v := userID.(type) {
	case uint:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return -1
	}
}
