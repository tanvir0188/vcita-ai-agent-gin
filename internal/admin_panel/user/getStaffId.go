package user

import (
	"fmt"
	"strings"

	"github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/utils"
	"github.com/tanvir0188/vcita-ai-agent/internal/types"
)

func GetStaffByEmail(email string) (*types.Staff, error) {
	staffs, err := utils.GetStaffList()
	if err != nil {
		return nil, err
	}

	for _, staff := range staffs {
		if strings.EqualFold(staff.Email, email) {
			return &staff, nil
		}
	}

	return nil, fmt.Errorf("staff not found")
}
