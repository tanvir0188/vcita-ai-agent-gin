package appointment

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/tanvir0188/vcita-ai-agent/internal/admin_panel/utils"
	"github.com/tanvir0188/vcita-ai-agent/internal/config"
	medicationreminder "github.com/tanvir0188/vcita-ai-agent/internal/medicationReminder"
)

func GetServices() ([]Service, error) {
	url := fmt.Sprintf(
		"https://api.vcita.biz/platform/v1/services?business_id=%s",
		config.Envs.BusinessUid,
	)

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + config.Envs.VcitaBusinessToken,
	}

	var response ServicesResponse

	if err := medicationreminder.DoRequest(http.MethodGet, url, headers, &response); err != nil {
		return nil, err
	}

	return response.Data.Services, nil
}

func GetAvailabilitySlots(
	serviceID string,
	startDate string,
	endDate string,
) (*AvailabilityResponse, error) {

	staffIDs, err := utils.GetStaffIDs()
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("id", serviceID)
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)

	for _, id := range staffIDs {
		params.Add("staff_ids[]", id)
	}

	endpoint := "https://api.vcita.biz/platform/v1/services/availability?" + params.Encode()

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + config.Envs.VcitaBusinessToken,
	}

	var response AvailabilityResponse
	if err := medicationreminder.DoRequest(http.MethodGet, endpoint, headers, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func BookAppointment(token string, req BookingRequest) (*BookingDetails, error) {
	url := "https://api.vcita.biz/business/scheduling/v1/bookings"

	headers := map[string]string{
		"accept":        "application/json",
		"authorization": "Bearer " + token,
	}

	var apiResponse bookingApiResponse

	// Call your existing helper function
	err := medicationreminder.DoRequestWithBody("POST", url, headers, req, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to book appointment: %w", err)
	}

	// Flatten out the important fields into the simple return struct
	b := apiResponse.Data.Booking
	return &BookingDetails{
		Title:              b.Title,
		StartTime:          b.StartTime,
		Duration:           b.Duration,
		ClientID:           b.ClientID,
		StaffID:            b.StaffID,
		ConversationID:     b.ConversationID,
		InteractionDetails: "01853958635",
	}, nil
}
