package user_queue

import (
	userModel "anonymous_chat/internal/models/user"
	"sync"
)

var (
	instance *InMemoryUserQueueStore
	once     sync.Once
)

type InMemoryUserQueueStore struct {
	list map[string]*userModel.User
	mu   sync.Mutex
}

func GetInstance() *InMemoryUserQueueStore {
	once.Do(func() {
		instance = newInMemoryUserQueueStore()
	})

	return instance
}

func newInMemoryUserQueueStore() *InMemoryUserQueueStore {
	return &InMemoryUserQueueStore{
		list: make(map[string]*userModel.User),
	}
}

func (s *InMemoryUserQueueStore) Add(user *userModel.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list[user.Hash] = user
}

func (s *InMemoryUserQueueStore) GetAll() map[string]*userModel.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.list
}

func (s *InMemoryUserQueueStore) Get(hash string) (*userModel.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, found := s.list[hash]

	return user, found
}

func (s *InMemoryUserQueueStore) Delete(hash string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.list, hash)
}
