package dto

type ContactList struct {
	Contacts []ContactForList `json:"contacts"`
}

type ContactForList struct {
	RoomID             int    `json:"room_id"`
	RecipientName      string `json:"recipient_name"`
	LastMessageContent string `json:"last_message_content"`
	LastMessageTime    string `json:"last_message_time"`
}
