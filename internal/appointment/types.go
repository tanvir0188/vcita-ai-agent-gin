package appointment

type ServicesResponse struct {
	Status string       `json:"status"`
	Data   ServicesData `json:"data"`
}

type ServicesData struct {
	Services []Service `json:"services"`
}

type Service struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	ProvidersStaff []string `json:"providers_staff"`
}

type BookingDetails struct {
	Title              string
	StartTime          string
	Duration           int
	ClientID           string
	StaffID            string
	ConversationID     string
	InteractionDetails string
}

// BookingRequest matches the API's required JSON payload
type BookingRequest struct {
	BusinessID         string `json:"business_id"`
	ClientID           string `json:"client_id"`
	MatterUID          string `json:"matter_uid"`
	ServiceID          string `json:"service_id"`
	StaffID            string `json:"staff_id"`
	StartTime          string `json:"start_time"`
	InteractionDetails string `json:"interaction_details"`
}

// bookingApiResponse handles the full nested structure from the server
type bookingApiResponse struct {
	Data struct {
		Booking struct {
			Title              string `json:"title"`
			StartTime          string `json:"start_time"`
			Duration           int    `json:"duration"`
			ClientID           string `json:"client_id"`
			StaffID            string `json:"staff_id"`
			ConversationID     string `json:"conversation_id"`
			InteractionDetails string `json:"interaction_details"`
		} `json:"booking"`
	} `json:"data"`
}
type Slot struct {
	StartTime string
	StaffID   string
}

// internal/appointment/types.go

// internal/appointment/types.go

type AvailabilityResponse struct {
	Status string           `json:"status"`
	Data   AvailabilityData `json:"data"`
}

type AvailabilityData struct {
	Availabilities map[string][]AvailabilitySlot `json:"availabilities"`
}

type AvailabilitySlot struct {
	StaffID         string  `json:"staff_id"`
	SpotsTotal      int     `json:"spots_total"`
	SpotsOpen       int     `json:"spots_open"`
	StartTime       string  `json:"start_time"`
	DurationMinutes int     `json:"duration_minutes"`
	EventInstanceID *string `json:"event_instance_id"`
}
