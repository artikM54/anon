package handler_queue

import (
	userModel "anonymous_chat/internal/models/user"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
	"math/rand"
	"time"
)

var userQueue = make(map[string]*userModel.User)

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
	userQueue[user.Hash] = user
}

func ExitUserWithinQueue(user *userModel.User) bool {
	_, found := userQueue[user.Hash]

	return found
}
	
func DeleteUserFromQueue(userHash string) {
	delete(userQueue, userHash)
}

func bindUser() {
	fmt.Println("There are two users")

	users := chooseRandomUsers(userQueue, 2)

	for _, user := range users {
		DeleteUserFromQueue(user.Hash)
	}

	c := chatService.NewChatService(users, &userQueue)

	go c.Start()
}

func chooseRandomUsers(m map[string]*userModel.User, n int) []*userModel.User {
	rand.Seed(time.Now().UnixNano())

	values := make([]*userModel.User, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}

	rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })

	if n > len(values) {
		n = len(values)
	}

	result := values[:n]

	return result
}
