package chat

import (
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
)

var userQueue []*userModel.User

func AddUserToQueue(user *userModel.User) {
	userQueue = append(userQueue, user)
}

func GetCountUsersIntoQueue() int {
	return len(userQueue)
}

func BindClients() {
	fmt.Println("There are two users")

	user1, user2, err := chooseRandomPair()
	if err != nil {
		log.Println("Error choosing random pair:", err)
		return
	}

	userQueue = removeClientFromSlice(userQueue, user1)
	userQueue = removeClientFromSlice(userQueue, user2)

	go HandleStreamMessages(user1, user2)
	go HandleStreamMessages(user2, user1)
}

func chooseRandomPair() (*userModel.User, *userModel.User, error) {
	idx1 := rand.Intn(len(userQueue))
	idx2 := rand.Intn(len(userQueue))

	for idx2 == idx1 {
		idx2 = rand.Intn(len(userQueue))
	}

	return userQueue[idx1], userQueue[idx2], nil
}

func removeClientFromSlice(slice []*userModel.User, user *userModel.User) []*userModel.User {
	for i, c := range slice {
		if c == user {
			return append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}

func HandleStreamMessages(user1 *userModel.User, user2 *userModel.User) {
	for {
		_, textMessage, err := user1.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client 1:", err)
			return
		}

		message := messageModel.NewMessage(
			"chat",
			string(textMessage),
		)

		if err := SendMessage(user2, message); err != nil {
			log.Println("Error sending message to client 2: ", err)
			return
		}
	}
}

func SendMessage(user *userModel.User, message *messageModel.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message to JSON: ", err)
	}

	return user.Conn.WriteMessage(websocket.TextMessage, data)
}
