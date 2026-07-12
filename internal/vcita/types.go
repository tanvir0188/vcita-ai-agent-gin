package vcita

// ── Message types ─────────────────────────────────────────────────────────────

type CreateMessageRequest struct {
	Message MessagePayload `json:"message"`
}

type MessagePayload struct {
	Direction string `json:"direction"`
	ClientID  string `json:"client_id"`
	Text      string `json:"text"`
}

type Message struct {
	UId             string `json:"uid"`
	Text            string `json:"text"`
	ConversationUID string `json:"conversation_uid"`

	Staff struct {
		Uid   string `json:"uid"`
		Email string `json:"email"`
	} `json:"staff"`

	WasRead   bool   `json:"was_read"`
	Direction string `json:"direction"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ── Smart Reply types ─────────────────────────────────────────────────────────

type SmartReplyRequest struct {
	MatterUID string `json:"matter_uid"`
	ClientUID string `json:"client_uid"`
}

type SmartReplyResponse struct {
	Success bool           `json:"success"`
	Data    SmartReplyData `json:"data"`
}

type SmartReplyData struct {
	UID       string            `json:"uid"`
	MatterUID string            `json:"matter_uid"`
	Payload   SmartReplyPayload `json:"payload"`
}

type SmartReplyPayload struct {
	EmailMessage string `json:"email_message"`
	SMSMessage   string `json:"sms_message"`
	FBMessage    string `json:"FB_message"`
	OtherMessage string `json:"other_message"`
}

// ── Client / Contact types ────────────────────────────────────────────────────

type GetClientDetailResponse struct {
	Status string           `json:"status"`
	Data   ClientDetailData `json:"data"`
}

type ClientDetailData struct {
	Client ClientInfo `json:"client"`
}

type ClientInfo struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobile_phone"`
}

// NotesResponse is used by the business notes endpoint.
type NotesResponse struct {
	Data struct {
		Notes []struct {
			Content string `json:"content"`
		} `json:"notes"`
	} `json:"data"`
}

// ── Staff types ───────────────────────────────────────────────────────────────

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
