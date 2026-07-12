package vcita

import (
	"fmt"
	"net/http"

	"github.com/tanvir0188/vcita-ai-agent/internal/config"
)

// GetStaffList fetches all staff members for the business from vcita.
func GetStaffList() ([]Staff, error) {
	businessUid := config.Envs.BusinessUid

	url := fmt.Sprintf(
		"https://api.vcita.biz/platform/v1/businesses/%s/staffs",
		businessUid,
	)

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + config.Envs.VcitaBusinessToken,
	}

	var response GetStaffResponse

	if err := DoRequest(http.MethodGet, url, headers, &response); err != nil {
		return nil, err
	}

	return response.Data.Staff, nil
}

// GetStaffIDs returns a list of active, non-deleted staff IDs.
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
