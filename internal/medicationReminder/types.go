package medicationreminder

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
