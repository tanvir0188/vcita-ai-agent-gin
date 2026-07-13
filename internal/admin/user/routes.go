package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tanvir0188/vcita-ai-agent/internal/admin/middleware"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/shared"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(rg *gin.RouterGroup, s *Store) {
	rg.GET("/login", ShowLoginPage)
	rg.POST("/login", s.HandleAdminLogin)
	rg.POST("/register", s.HandleRegister)

	protected := rg.Group("/")
	protected.Use(WithJWTAuth(s))

	protected.GET("/users", s.ListUsersPage)
	protected.PATCH("/profile", s.HandleProfile)
	protected.GET("/profile", s.HandleGetProfile)
	protected.PATCH("/change-password", s.HandleChangePassword)
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title": "Admin Login",
	})
}

func (s *Store) HandleAdminLogin(c *gin.Context) {
	var payload LoginUserPayload

	if err := shared.BindAndValidate(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user store.User

	err := s.db.
		Where("email = ?", payload.Email).
		First(&user).
		Error

	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusBadRequest,
			Message: "No user found with the given email",
			Err:     err,
		})
		c.Abort()
		return
	}

	if !user.IsActive {
		c.Error(&middleware.AppError{
			Status:  http.StatusForbidden,
			Message: "account is inactive",
			Err:     err,
		})
		c.Abort()
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
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to generate token",
			Err:     err,
		})
		c.Abort()
		return
	}

	refreshToken, err := CreateRefreshToken(secret, uint(user.ID))
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to generate refresh token",
			Err:     err,
		})
		c.Abort()
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
	var payload RegisterPayload

	if err := shared.BindAndValidate(c, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
	staff, err := GetStaffByEmail(payload.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to register in this admin panel",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(payload.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to hash password",
			Err:     err,
		})
		c.Abort()
		return
	}
	fullName := staff.DisplayName
	phoneNumber := staff.MobileNumber

	user := store.User{
		StaffUID:    staff.ID,
		FullName:    fullName,
		Email:       payload.Email,
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),

		IsActive:   true,
		IsVerified: true,
	}

	err = s.db.Create(&user).Error
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to create user",
			Err:     err,
		})
		c.Abort()
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

	err := s.db.Find(&users).Error
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to retrieve users",
			Err:     err,
		})
		c.Abort()
		return
	}

	var response []UserListResponse

	for _, u := range users {
		response = append(response, UserListResponse{
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

func (s *Store) HandleProfile(c *gin.Context) {
	var payload ProfilePayload

	if err := shared.BindAndValidate(c, &payload); err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusBadRequest,
			Message: "failed to parse request payload",
			Err:     err,
		})
		c.Abort()
		return
	}

	user, err := shared.GetUserFromRequest(c)
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusUnauthorized,
			Message: "failed to get user from request",
			Err:     err,
		})
		c.Abort()
		return
	}

	fmt.Println(user.Email)
	fmt.Println(user.IsAdmin)

	if err := s.db.First(&user, &user.ID).Error; err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusNotFound,
			Message: "user not found",
			Err:     err,
		})
		c.Abort()
		return
	}

	updates := map[string]interface{}{}

	if payload.FullName != nil {
		updates["full_name"] = *payload.FullName
	}
	if payload.Email != nil {
		updates["email"] = *payload.Email
	}
	if payload.PhoneNumber != nil {
		updates["phone_number"] = *payload.PhoneNumber
	}

	if len(updates) == 0 {
		c.Error(&middleware.AppError{
			Status:  http.StatusBadRequest,
			Message: "no fields to update",
		})
		c.Abort()
		return
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update profile",
			Err:     err,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data": gin.H{
			"full_name":    user.FullName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
		},
	})
}

func (s *Store) HandleGetProfile(c *gin.Context) {
	user, err := shared.GetUserFromRequest(c)
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusUnauthorized,
			Message: "failed to get user from request",
			Err:     err,
		})
		c.Abort()
		return
	}

	if err := s.db.First(&user, &user.ID).Error; err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusNotFound,
			Message: "user not found",
			Err:     err,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"full_name":    user.FullName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"is_admin":     user.IsAdmin,
			"id":           user.ID,
		},
	})
}

func (s *Store) HandleChangePassword(c *gin.Context) {
	var payload PasswordPayload

	if err := shared.BindAndValidate(c, &payload); err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusBadRequest,
			Message: "failed to parse request payload",
			Err:     err,
		})
		c.Abort()
		return
	}

	user, err := shared.GetUserFromRequest(c)
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusUnauthorized,
			Message: "failed to get user from request",
			Err:     err,
		})
		c.Abort()
		return
	}

	fmt.Println(user.Email)
	fmt.Println(user.IsAdmin)

	if err := s.db.First(&user, &user.ID).Error; err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusNotFound,
			Message: "user not found",
			Err:     err,
		})
		c.Abort()
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(payload.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to hash password",
			Err:     err,
		})
		c.Abort()
		return
	}

	user.Password = string(hashedPassword)

	if err := s.db.Save(&user).Error; err != nil {
		c.Error(&middleware.AppError{
			Status:  http.StatusInternalServerError,
			Message: "failed to update password",
			Err:     err,
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password changed successfully",
	})
}
