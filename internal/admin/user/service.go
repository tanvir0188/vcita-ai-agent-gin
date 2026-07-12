package user

import (
	"fmt"
	"strings"

	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
)

// GetStaffByEmail looks up a staff member by their email address.
func GetStaffByEmail(email string) (*vcita.Staff, error) {
	staffs, err := vcita.GetStaffList()
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
