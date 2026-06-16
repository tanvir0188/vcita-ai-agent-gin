package types

type User struct {
	StaffUID    string `json:"staff_uid"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"-"`
	IsActive    bool   `json:"is_active"`
	IsVerified  bool   `json:"is_verified"`
}
