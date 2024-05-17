package message

type MessagePayload struct {
	Text      string `json:"text"`
	Timestamp int64 `json:"timestamp"`
	UserHash  string `json:"userHash"`
	ChatHash  string `json:"chatHash"`
}
type Message struct {
	Category string          `json:"type"`
	Payload  *MessagePayload `json:"payload"`
}

func NewMessage(category string, text string, timestamp int64, userHash string, chatHash string) *Message {
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
