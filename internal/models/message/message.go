package message

import (
	"time"
)

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

func NewMessage(category string, text string, userHash string, chatHash string) *Message {
	now := time.Now()
	unixTimestamp := now.Unix()

	payload := &MessagePayload{
		Text:      text,
		Timestamp: unixTimestamp,
		UserHash:  userHash,
		ChatHash:  chatHash,
	}

	return &Message{
		Category: category,
		Payload:  payload,
	}
}
