package chat_list

import (
	chatModel "anonymous_chat/internal/models/chat"
	"sync"
)

var (
	instance *InMemoryChatListStore
	once     sync.Once
)

type InMemoryChatListStore struct {
	list map[string]*chatModel.Chat
	mu   sync.Mutex
}

func GetInstance() *InMemoryChatListStore {
	once.Do(func() {
		instance = newInMemoryChatListStore()
	})

	return instance
}

func newInMemoryChatListStore() *InMemoryChatListStore {
	return &InMemoryChatListStore{
		list: make(map[string]*chatModel.Chat),
	}
}

func (s *InMemoryChatListStore) Add(chat *chatModel.Chat) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list[chat.Hash] = chat
}

func (s *InMemoryChatListStore) Get(hash string) (*chatModel.Chat, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	chat, found := s.list[hash]

	return chat, found
}

func (s *InMemoryChatListStore) Delete(hash string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.list, hash)
}
