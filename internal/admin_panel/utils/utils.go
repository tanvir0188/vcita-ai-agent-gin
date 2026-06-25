package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	medicationreminder "github.com/tanvir0188/vcita-ai-agent/internal/medicationReminder"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/types"
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

func GetStaffList() ([]types.Staff, error) {
	businessUid := config.Envs.BusinessUid

	url := fmt.Sprintf(
		"https://api.vcita.biz/platform/v1/businesses/%s/staffs",
		businessUid,
	)

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + config.Envs.VcitaBusinessToken,
	}

	var response types.GetStaffResponse

	if err := medicationreminder.DoRequest(http.MethodGet, url, headers, &response); err != nil {
		return nil, err
	}

	return response.Data.Staff, nil
}

func GetStaffIDs() ([]string, error) {
	staffs, err := GetStaffList()
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(staffs))

	for _, staff := range staffs {
		if staff.Active && !staff.Deleted {
			ids = append(ids, staff.ID)
		}
	}

	return ids, nil
}
