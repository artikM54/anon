package message

type MessagePayload struct {
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
	UserHash  string `json:"userHash"`
	ChatHash  string `json:"chatHash"`
}
type Message struct {
	Category string          `json:"type"`
	Payload  *MessagePayload `json:"payload"`
}

func NewMessage(category string, text string, timestamp string, userHash string, chatHash string) *Message {
	payload := &MessagePayload{
		Text:      text,
		Timestamp: timestamp,
		UserHash:  userHash,
		ChatHash:  chatHash,
	}

	return &Message{
		Category: category,
		Payload:  payload,
	}
}
