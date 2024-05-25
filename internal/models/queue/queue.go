package queue

import (
	userModel "anonymous_chat/internal/models/user"
	"math/rand"
	"time"
	"sync"
)

type UserQueue struct {
	users map[string]*userModel.User
	mu    sync.Mutex
}

func NewUserQueue() UserQueue {
	return UserQueue{
		users: make(map[string]*userModel.User),
	}
}

func (q *UserQueue) GetCountUsersIntoQueue() int {
	q.mu.Lock()
	len := len(q.users)
	q.mu.Unlock()

	return len
}

func (q *UserQueue) AddUserToQueue(user *userModel.User) {
	q.mu.Lock()
	q.users[user.Hash] = user
	q.mu.Unlock()
}

func (q *UserQueue) ExitUserWithinQueue(userHash string) bool {
	q.mu.Lock()
	_, found := q.users[userHash]
	q.mu.Unlock()

	return found
}

func (q *UserQueue) DeleteUserFromQueue(userHash string) {
	q.mu.Lock()
	delete(q.users, userHash)
	q.mu.Unlock()
}

func (q *UserQueue) ChooseRandomUsers(n int) []*userModel.User {
	rand.Seed(time.Now().UnixNano())

	values := make([]*userModel.User, 0, len(q.users))
	for _, value := range q.users {
		values = append(values, value)
	}

	rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })

	if n > len(values) {
		n = len(values)
	}

	result := values[:n]

	for _, user := range result {
		q.DeleteUserFromQueue(user.Hash)
	}

	return result
}