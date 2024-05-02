package message

type Message struct {
	Category string `json:"category"`
	Text     string `json:"text"`
}

func NewMessage(category string, text string) *Message {
	return &Message{
		Category: category,
		Text:     text,
	}
}
