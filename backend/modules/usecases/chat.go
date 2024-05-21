package usecases

import (
	"github.com/bukharney/giga-chat/modules/entities"
)

type ChatUsecase struct {
	ChatRepo entities.ChatRepository
}

func NewChatUsecases(chatRepo entities.ChatRepository) entities.ChatUsecase {
	return &ChatUsecase{ChatRepo: chatRepo}
}

func (c *ChatUsecase) CreateChatRoom(req *entities.ChatRoom) error {
	return nil
}

func (c *ChatUsecase) GetChatMessages(roomId int) ([]entities.ChatMessage, error) {
	message, err := c.ChatRepo.GetChatMessages(roomId)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (c *ChatUsecase) GetChatRoom(userId int, roomId int) error {
	err := c.ChatRepo.GetChatRoom(userId, roomId)
	if err != nil {
		return err
	}

	return nil
}
