package handler_queue

import (
	userModel "anonymous_chat/internal/models/user"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
	"log"
	"math/rand"
)

var userQueue []*userModel.User

func MustLoad() {
	go start()
}

func start() {
	for {
		if getCountUsersIntoQueue() >= 2 {
			bindUser()
		}
	}
}

func getCountUsersIntoQueue() int {
	return len(userQueue)
}

func AddUserToQueue(user *userModel.User) {
	userQueue = append(userQueue, user)
}

func bindUser() {
	fmt.Println("There are two users")

	users, err := chooseRandomUsers()
	if err != nil {
		log.Println("Error choosing random pair:", err)
		return
	}

	for _, user := range users {
		userQueue = removeClientFromSlice(userQueue, user)
	}

	c := chatService.NewChatService(users, &userQueue)

	go c.Start()
}

func chooseRandomUsers() ([]*userModel.User, error) {
	result := make([]*userModel.User, 0, 5)

	idx1 := rand.Intn(len(userQueue))
	idx2 := rand.Intn(len(userQueue))

	for idx2 == idx1 {
		idx2 = rand.Intn(len(userQueue))
	}
	result = append(result, userQueue[idx1])
	result = append(result, userQueue[idx2])

	return result, nil
}

func removeClientFromSlice(slice []*userModel.User, user *userModel.User) []*userModel.User {
	for i, c := range slice {
		if c == user {
			return append(slice[:i], slice[i+1:]...)
		}
	}

	return slice
}
