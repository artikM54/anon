package user_queue

import (
	"anonymous_chat/internal/repositories/user_queue"
	chatService "anonymous_chat/internal/services/chat"
	"fmt"
)

type UserQueueService struct {
	repository *user_queue.UserQueueRepository
}

func NewUserQueueService() *UserQueueService {
	return &UserQueueService{
		repository: user_queue.NewUserQueueRepository(),
	}
}

func (s *UserQueueService) Start() {
	go s.handler()
}

func (s *UserQueueService) handler() {
	for {
		if s.repository.GetCountUsers() >= 2 {
			go s.bindUser()
		}
	}
}

func (s *UserQueueService) bindUser() {
	fmt.Println("There are two users")

	users := s.repository.GetRandomUsers(2)

	c := chatService.NewChatService(users)

	c.Start()
	fmt.Println("There are two users END")
}
