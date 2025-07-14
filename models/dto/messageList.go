package dto

type MessageList struct {
	Type     string            `json:"type"`
	ChatName string            `json:"chatName"`
	Messages []MessageResponse `json:"messages"`
}

type MessageResponse struct {
	MessageId int    `json:"messageId"`
	Content   string `json:"content"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
}
