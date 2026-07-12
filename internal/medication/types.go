package medication

type ClientNotesResponse struct {
	Success bool            `json:"success"`
	Data    ClientNotesData `json:"data"`
}

type ClientNotesData struct {
	ClientNotes []ClientNote `json:"client_notes"`
}

type ClientNote struct {
	UID       string  `json:"uid"`
	MatterUID string  `json:"matter_uid"`
	Title     *string `json:"title"`
	Content   string  `json:"content"`
	StaffUID  string  `json:"staff_uid"`
}

type TrimmedNote struct {
	NoteID  string `json:"note_id"`
	Content string `json:"content"`
}

type TrimmedNotesResponse struct {
	MatterUID    string        `json:"matter_uid"`
	TrimmedNotes []TrimmedNote `json:"trimmed_notes"`
}

type ClientNoteResponse struct {
	Success bool `json:"success"`
	Data    struct {
		UID     string `json:"uid"`
		Content string `json:"content"`
	} `json:"data"`
}

type Client struct {
	StaffIds  int    `json:"staff_ids"`
	Notes     string `json:"notes"`
	MatterUid string `json:"matter_uid"`
	UID       string `json:"uid"`
}

type FilteredClient struct {
	StaffIds  int    `json:"staff_ids"`
	Notes     string `json:"notes"`
	MatterUid string `json:"matter_uid"`
}

type ClientListResponse struct {
	Counts struct {
		Client int `json:"client"`
	} `json:"counts"`

	TopHits struct {
		Client []Client `json:"client"`
	} `json:"top_hits"`
}
