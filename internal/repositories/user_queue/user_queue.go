package user_queue

import (
	"anonymous_chat/internal/memory_storages/user_queue"
	userModel "anonymous_chat/internal/models/user"
	"fmt"
	"math/rand"
	"time"
)

type UserQueueRepository struct {
	store *user_queue.InMemoryUserQueueStore
}

func NewUserQueueRepository() *UserQueueRepository {
	return &UserQueueRepository{
		store: user_queue.GetInstance(),
	}
}

func (r *UserQueueRepository) Add(user *userModel.User) {
	r.store.Add(user)
}

func (r *UserQueueRepository) Get(hash string) *userModel.User {
	user, found := r.store.Get(hash)
	if !found {
		fmt.Println("Repository queue could not get a user")
	}

	return user
}

func (r *UserQueueRepository) Delete(hash string) {
	r.store.Delete(hash)
}

func (r *UserQueueRepository) Exist(hash string) bool {
	_, found := r.store.Get(hash)

	return found
}

func (r *UserQueueRepository) GetCountUsers() int {
	return len(r.store.GetAll())
}

func (r *UserQueueRepository) GetRandomUsers(n int) map[string]*userModel.User {
	rand.Seed(time.Now().UnixNano())
	users := r.store.GetAll()

	values := make([]*userModel.User, 0, len(users))
	for _, value := range users {
		values = append(values, value)
	}

	rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })

	if n > len(values) {
		n = len(values)
	}

	randomUsers := values[:n]
	result := make(map[string]*userModel.User)

	for _, user := range randomUsers {
		result[user.Hash] = user
		r.store.Delete(user.Hash)
	}

	return result
}
