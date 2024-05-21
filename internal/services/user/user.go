package user

import (
	"anonymous_chat/internal/handler_queue"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"

	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type UserService struct {
	user *userModel.User
}

func NewUserService(user *userModel.User) *UserService {
	return &UserService{
		user: user,
	}
}

func (u *UserService) Start() {
	go u.listeningMessages()
	go u.spendingMessages()
}
func (u *UserService) listeningMessages() {
	for {
		fmt.Println("HANDLER user messages ", u.user.Hash)

		data, err := u.user.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message for user %s: %v\n", u.user.Hash, err)
			return
		}

		var message messageModel.Message
		err = json.Unmarshal(data, &message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				fmt.Printf("WebSocket closed unexpectedly for user %s: %v\n", u.user.Hash, err)
			} else {
				fmt.Printf("Error reading WebSocket message for user %s: %v\n", u.user.Hash, err)
			}
		}

		now := time.Now()
		unixTimestamp := now.Unix()

		message.Payload.Timestamp = unixTimestamp
		message.Payload.UserHash = u.user.Hash
		fmt.Println("HANDLER user messages message chat Hash: ", message.Payload.ChatHash)
		switch message.Category {
		case "FRONT:START_QUEUE":
			fmt.Printf("HANDLE COMMANDS FRONT:START_QUEUE for user %s\n", u.user.Hash)

			if !handler_queue.ExitUserWithinQueue(u.user.Hash) {
				handler_queue.AddUserToQueue(u.user)
			}
		case "FRONT:EXIT_QUEUE":
			fmt.Printf("HANDLE COMMANDS FRONT:EXIT_QUEUE for user %s\n", u.user.Hash)
			handler_queue.DeleteUserFromQueue(u.user.Hash)
		case "CHAT":
			fmt.Printf("HANDLE COMMANDS CHAT for user %s\n", u.user.Hash)

			if message.Payload.ChatHash != "" {
				u.user.PutToInChat(&message)
			}
		case "FRONT:CHAT_EXIT":
			fmt.Printf("HANDLE COMMANDS FRONT:CHAT_EXIT for user %s\n", u.user.Hash)

			if message.Payload.ChatHash != "" {
				message.Category = "CHAT_EXIT"
				u.user.PutToInChat(&message)
				u.user.CloseChat(message.Payload.ChatHash)
			}

		default:
			fmt.Printf("Undefined command for user %s: %s\n", u.user.Hash, message.Category)
		}
	}
}

func (u *UserService) spendingMessages() {
	for {
		message, opened := u.user.GetFromOutChannel()
		if !opened {
			fmt.Println("out channel is closed")
			return
		}

		data, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error encoding message to JSON: ", err)
		}

		fmt.Println("spendingMessages() user: ", u.user.Hash, "message data: ", message)

		u.user.Conn.WriteMessage(websocket.TextMessage, data)
	}
}
