package message

const (
	ConnectCategory         = "CONNECT"
	TokenCategory           = "TOKEN"
	StartCategory           = "CHAT_START"
	ChatCategory            = "CHAT"
	ExitCategory            = "CHAT_EXIT"
	FrontStartQueueCategory = "FRONT:START_QUEUE"
	FrontExitQueueCategory  = "FRONT:EXIT_QUEUE"
	FrontChatExitCategory   = "FRONT:CHAT_EXIT"
	FrontGetTokenCategory   = "FRONT:GET_TOKEN"
	FrontGiveTokenCategory  = "FRONT:GIVE_TOKEN"
)

type MessagePayload struct {
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
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
