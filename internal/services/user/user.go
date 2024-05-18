package user

import (
	"anonymous_chat/internal/handler_queue"
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	hashUtil "anonymous_chat/internal/utils/hash"

	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type UserService struct {
	user *userModel.User
}

func NewUser(conn *websocket.Conn, hash string) *userModel.User {
	if hash == "" {
		hash = hashUtil.CreateUniqueModelHash(userModel.RedisList)
	}

	return &userModel.User{
		Hash: hash,
		Conn: conn,
	}
}

func NewUserService(user *userModel.User) *UserService {
	return &UserService{
		user: user,
	}
}

func (c *UserService) HandleUsersCommand() {
	for {
		fmt.Println("HANDLER users command start service ", c.user.Hash)

		data, err := c.user.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message for user %s: %v\n", c.user.Hash, err)
			return
		}

		var message messageModel.Message
		err = json.Unmarshal(data, &message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				fmt.Printf("WebSocket closed unexpectedly for user %s: %v\n", c.user.Hash, err)
			} else {
				fmt.Printf("Error reading WebSocket message for user %s: %v\n", c.user.Hash, err)
			}
		}

		switch message.Category {
		case "FRONT:START_QUEUE":
			fmt.Printf("HANDLE COMMANDS FRONT:START_QUEUE for user %s\n", c.user.Hash)
			c.user.ChannelChat = make(chan *messageModel.Message, 30)
			handler_queue.AddUserToQueue(c.user)
		case "FRONT:EXIT_QUEUE":
			fmt.Printf("HANDLE COMMANDS FRONT:EXIT_QUEUE for user %s\n", c.user.Hash)
			handler_queue.DeleteUserFromQueue(c.user.Hash)
		case "CHAT":
			fmt.Printf("HANDLE COMMANDS CHAT for user %s\n", c.user.Hash)
			c.user.ChannelChat <- &message
		case "FRONT:CHAT_EXIT":
			fmt.Printf("HANDLE COMMANDS FRONT:CHAT_EXIT for user %s\n", c.user.Hash)
			c.user.ChannelChat <- &message
			close(c.user.ChannelChat)
		default:
			fmt.Printf("Undefined command for user %s: %s\n", c.user.Hash, message.Category)
		}
	}
}
