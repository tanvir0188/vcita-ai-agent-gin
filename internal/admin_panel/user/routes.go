package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/utils"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(rg *gin.RouterGroup, s *Store) {
	rg.GET("/login", ShowLoginPage)
	rg.POST("/login", s.HandleAdminLogin)
	rg.POST("/register", s.HandleRegister)

	protected := rg.Group("/")
	protected.Use(WithJWTAuth(*s))

	protected.GET("/users", s.ListUsersPage)

	//rg.GET("/users/create", ShowCreateUserPage)

	// rg.POST("/users", CreateUser)

	// rg.GET("/users/:id", ShowUserPage)
	// rg.GET("/users/:id/edit", ShowEditUserPage)

	// rg.POST("/users/:id", UpdateUser)
	// rg.POST("/users/:id/delete", DeleteUser)
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title": "Admin Login",
	})
}

func (s *Store) HandleAdminLogin(c *gin.Context) {
	var payload types.LoginUserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid payload: %v", validationErrors),
		})
		return
	}

	var user store.User

	err := s.db.
		Where("email = ?", payload.Email).
		First(&user).
		Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No user found with the given email",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "account is inactive",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(payload.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	secret := []byte(config.Envs.JWTSecret)

	accessToken, err := CreateAccessToken(secret, uint(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := CreateRefreshToken(secret, uint(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"email":     user.Email,
			"full_name": user.FullName,
			"is_admin":  user.IsAdmin,
		},
	})
}

func (s *Store) HandleRegister(c *gin.Context) {
	var payload types.RegisterPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Println(err)
		return
	}

	if err := utils.Validate.Struct(&payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid payload: %v", validationErrors),
		})
		return
	}

	var existingUser store.User

	err := s.db.
		Where("email = ?", payload.Email).
		First(&existingUser).
		Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(payload.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := store.User{
		StaffUID:    payload.StaffUID,
		FullName:    payload.FullName,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Password:    string(hashedPassword),

		IsActive:   true,
		IsVerified: true,
	}

	err = s.db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user": gin.H{
			"id":           user.ID,
			"staff_uid":    user.StaffUID,
			"full_name":    user.FullName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"is_active":    user.IsActive,
			"is_verified":  user.IsVerified,
		},
	})
}

func (s *Store) ListUsersPage(c *gin.Context) {
	var users []store.User

	log.Println("ListUsersPage API called")

	err := s.db.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve users",
		})
		return
	}
	var response []types.UserListResponse
	for _, u := range users {
		response = append(response, types.UserListResponse{
			ID:          u.ID,
			StaffUID:    u.StaffUID,
			FullName:    u.FullName,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			IsActive:    u.IsActive,
			IsVerified:  u.IsVerified,
			IsAdmin:     u.IsAdmin,
			CreatedAt:   u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
