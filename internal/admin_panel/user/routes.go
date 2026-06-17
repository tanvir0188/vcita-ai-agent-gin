package user

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutes(rg *gin.RouterGroup, s *Store) {
	rg.GET("/login", ShowLoginPage)
	rg.POST("/login", s.HandleAdminLogin)

	rg.GET("/users", s.ListUsersPage)

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
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user store.User

	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"title": "Admin Login",
			"error": "Invalid credentials",
		})
		log.Printf("admin login failed for email %s: %v", email, err)
		return
	}

	if !user.IsActive {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"title": "Admin Login",
			"error": "Account is inactive",
		})
		return
	}

	if !user.IsAdmin {
		c.HTML(http.StatusForbidden, "login.tmpl", gin.H{
			"title": "Admin Login",
			"error": "Access denied",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"title": "Admin Login",
			"error": "Invalid credentials",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("admin_id", user.ID)
	session.Set("admin_logged_in", true)
	session.Save()

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (s *Store) ListUsersPage(c *gin.Context) {
	var users []store.User

	err := s.db.Find(&users).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "base.tmpl", gin.H{
			"title": "User List",
			"error": "Failed to retrieve users",
		})
		return
	}

	c.HTML(http.StatusOK, "base.tmpl", gin.H{
		"title": "User List",
		"users": users,
	})
}
