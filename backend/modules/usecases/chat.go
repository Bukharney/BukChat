package usecases

import (
	"github.com/bukharney/giga-chat/modules/entities"
)

type ChatUsecase struct {
	ChatRepo entities.ChatRepository
}

func NewChatUsecase(chatRepo entities.ChatRepository) entities.ChatUsecase {
	return &ChatUsecase{ChatRepo: chatRepo}
}

func (c *ChatUsecase) CreateChatRoom(req *entities.ChatRoom) error {
	return nil
}
