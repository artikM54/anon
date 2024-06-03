package return_chat

import (
	messageModel "anonymous_chat/internal/models/message"
	userModel "anonymous_chat/internal/models/user"
	chatRepository "anonymous_chat/internal/repositories/chat"
	"anonymous_chat/internal/repositories/chat_list"
	"encoding/json"
	"fmt"
)

type ReturnToChatService struct {
	chatHash                     string
	user                         *userModel.User
	chatMessageHistoryRepository *chatRepository.ChatRepository
	chatListRepository           *chat_list.ChatListRepository
}

func NewReturnToChatService(chatHash string, user *userModel.User) *ReturnToChatService {

	return &ReturnToChatService{
		chatHash:                     chatHash,
		user:                         user,
		chatMessageHistoryRepository: chatRepository.NewChatRepository(chatHash),
		chatListRepository:           chat_list.NewChatListRepository(),
	}
}

func (s *ReturnToChatService) GetHistoryMessages() {
	messages := s.chatMessageHistoryRepository.GetMessages()

	data, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("AddMessage Redis Error encoding message to JSON: ", err)
	}

	message := messageModel.NewMessage(
		"HISTORY",
		string(data),
		s.user.Hash,
		s.chatHash,
	)

	s.user.PutIntoChannel(message)
}

func (s *ReturnToChatService) ReturnToChat() {
	s.chatListRepository.AddUserToChat(s.chatHash, s.user)
}
