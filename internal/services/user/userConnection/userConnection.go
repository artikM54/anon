package userConnection

import (
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	"anonymous_chat/internal/repositories/chat_list"
	"anonymous_chat/internal/repositories/user_queue"
	hashUtil "anonymous_chat/internal/utils/hash"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type userConnectionService struct {
	conn                *websocket.Conn
	user                *userModel.User
	chatListRepository  *chat_list.ChatListRepository
	userQueueRepository *user_queue.UserQueueRepository
}

func NewUserConnectionService(conn *websocket.Conn) *userConnectionService {
	user := userModel.NewUser()

	return &userConnectionService{
		conn:                conn,
		user:                user,
		chatListRepository:  chat_list.NewChatListRepository(),
		userQueueRepository: user_queue.NewUserQueueRepository(),
	}
}

func (u *userConnectionService) Start() {
	go u.listening()
	go u.sending()
}

func (u *userConnectionService) listening() {
	for {
		_, data, err := u.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		message := u.prepearMessage(data)

		u.handleMessage(message)
	}
}

func (u *userConnectionService) sending() {
	for {
		message, opened := u.user.GetFromChannel()
		if !opened {
			fmt.Println("out channel is closed")
			return
		}

		data, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error encoding message to JSON: ", err)
		}

		fmt.Println("spendingMessages() user: ", u.user.Hash, "message data: ", message)

		u.conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (u *userConnectionService) prepearMessage(data []byte) *messageModel.Message {
	var message messageModel.Message

	err := json.Unmarshal(data, &message)
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

	fmt.Println("handler conn message", message)

	return &message
}

func (u *userConnectionService) handleMessage(message *messageModel.Message) {
	fmt.Println("handleMessage: category - ", message.Category)
	switch message.Category {
	case messageModel.FrontGetTokenCategory:
		u.handleCaseFrontGetTokenCategory()

	case messageModel.FrontGiveTokenCategory:
		u.handleCaseFrontGiveTokenCategory(message)

	case messageModel.FrontStartQueueCategory:
		u.handleCaseFrontStartQueueCategory()

	case messageModel.FrontExitQueueCategory:
		u.handleCaseFrontExitQueueCategory()

	case messageModel.ChatCategory:
		u.handleCaseChatCategory(message)

	case messageModel.FrontChatExitCategory:
		u.handleCaseFrontChatExitCategory(message)

	default:
		fmt.Println("Undefined command")
	}
}

func (u *userConnectionService) handleCaseFrontGetTokenCategory() {
	if u.user.Hash != "" {
		fmt.Println("Already user have a hash")
		return
	}

	hash := hashUtil.CreateUniqueModelHash(userModel.RedisList)
	u.user.SetToken(hash)

	message := messageModel.NewMessage(
		messageModel.TokenCategory,
		hash,
		"",
		"",
	)

	u.user.PutIntoChannel(message)
}

func (u *userConnectionService) handleCaseFrontGiveTokenCategory(message *messageModel.Message) {
	u.user.SetToken(message.Payload.Text)
}

func (u *userConnectionService) handleCaseFrontStartQueueCategory() {
	fmt.Printf("HANDLE COMMANDS FRONT:START_QUEUE for user %s\n", u.user.Hash)

	if u.userQueueRepository.Exist(u.user.Hash) {
		fmt.Printf("HANDLE COMMANDS FRONT:START_QUEUE THERE IS USER IN QUEUE %s\n", u.user.Hash)
		return
	}

	u.userQueueRepository.Add(u.user)
}

func (u *userConnectionService) handleCaseFrontExitQueueCategory() {
	fmt.Printf("HANDLE COMMANDS FRONT:QUEUE_EXIT for user %s\n", u.user.Hash)

	if !u.userQueueRepository.Exist(u.user.Hash) {
		fmt.Printf("HANDLE COMMANDS FRONT:QUEUE_EXIT there is not in queue for user %s\n", u.user.Hash)
		return
	}

	u.userQueueRepository.Delete(u.user.Hash)
}

func (u *userConnectionService) handleCaseChatCategory(message *messageModel.Message) {
	fmt.Printf("HANDLE COMMANDS CHAT for user %s\n", u.user.Hash)

	if message.Payload.ChatHash == "" {
		fmt.Printf("HANDLE COMMANDS CHAT HASH is empty; user %s\n", u.user.Hash)
		return
	}

	if !u.chatListRepository.ExistChat(message.Payload.ChatHash) {
		fmt.Printf("HANDLE COMMANDS CHAT HASH is not exist; user %s\n", u.user.Hash)
		return
	}

	u.chatListRepository.PutMessageIntoChat(message)
}

func (u *userConnectionService) handleCaseFrontChatExitCategory(message *messageModel.Message) {
	fmt.Printf("HANDLE COMMANDS FRONT:CHAT_EXIT for user %s\n", u.user.Hash)

	if message.Payload.ChatHash == "" {
		fmt.Printf("HANDLE COMMANDS FRONT:CHAT_EXIT HASH is empty; user %s\n", u.user.Hash)
		return
	}

	if !u.chatListRepository.ExistChat(message.Payload.ChatHash) {
		fmt.Printf("HANDLE COMMANDS FRONT:CHAT_EXIT HASH is not exist; user %s\n", u.user.Hash)
		return
	}

	message.Category = messageModel.ExitCategory

	u.chatListRepository.PutMessageIntoChat(message)
	u.chatListRepository.ExitUserFromChat(message.Payload.ChatHash, message.Payload.UserHash)
}
