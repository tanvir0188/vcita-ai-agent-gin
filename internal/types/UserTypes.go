package types

import "time"

type UserListResponse struct {
	ID          uint      `json:"id"`
	StaffUID    string    `json:"staff_uid"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active"`
	IsVerified  bool      `json:"is_verified"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type ProfilePayload struct {
	FullName    *string `json:"full_name" validate:"omitempty"`
	Email       *string `json:"email" validate:"omitempty,email"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty"`
}

type PasswordPayload struct {
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type GetStaffResponse struct {
	Status string       `json:"status"`
	Data   GetStaffData `json:"data"`
}

type GetStaffData struct {
	Staff []Staff `json:"staff"`
}

type Staff struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	Active       bool   `json:"active"`
	Deleted      bool   `json:"deleted"`
	MobileNumber string `json:"mobile_number"`
}
