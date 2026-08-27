package usecases

import (
	"context"

	"github.com/bukharney/bukchat/modules/entities"
)

type ChatUsecase struct {
	ChatRepo entities.ChatRepository
}

func NewChatUsecases(chatRepo entities.ChatRepository) entities.ChatUsecase {
	return &ChatUsecase{ChatRepo: chatRepo}
}

func (c *ChatUsecase) CreateChatRoom(ctx context.Context, req *entities.ChatRoom) error {
	_, err := c.ChatRepo.CreateChatRoom(ctx, req)
	return err
}

func (c *ChatUsecase) GetChatMessages(ctx context.Context, roomId int) ([]entities.ChatMessage, error) {
	message, err := c.ChatRepo.GetChatMessages(ctx, roomId)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (c *ChatUsecase) GetChatRoom(ctx context.Context, userId int, roomId int) error {
	err := c.ChatRepo.GetChatRoom(ctx, userId, roomId)
	if err != nil {
		return err
	}

	return nil
}
