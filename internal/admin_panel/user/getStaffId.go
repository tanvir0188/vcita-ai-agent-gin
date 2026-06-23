package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	"github.com/tanvir0188/vcita-ai-agent/internal/types"
)

func GetStaffByEmail(email string) (*types.Staff, error) {
	accessToken := "Bearer " + config.Envs.VcitaBusinessToken
	businessUid := config.Envs.BusinessUid

	url := "https://api.vcita.biz/platform/v1/businesses/" + businessUid + "/staffs"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response types.GetStaffResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	for _, staff := range response.Data.Staff {
		if strings.EqualFold(staff.Email, email) {
			return &staff, nil
		}
	}

	return nil, fmt.Errorf("staff not found")
}
