package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/jmoiron/sqlx"
)

type ChatRepo struct {
	Db *sqlx.DB
}

func NewChatRepo(db *sqlx.DB) entities.ChatRepository {
	return &ChatRepo{Db: db}
}

func (c *ChatRepo) CreateChatRoom(ctx context.Context, req *entities.ChatRoom) (int, error) {
	query := `INSERT INTO rooms (name) VALUES ($1) RETURNING id`
	var roomId int
	err := c.Db.QueryRowContext(ctx, query, req.Name).Scan(&roomId)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat room: %w", err)
	}

	return roomId, nil
}

func (c *ChatRepo) GetChatRoom(ctx context.Context, userId int, roomId int) error {
	query := `SELECT user_id, room_id FROM users_rooms
	WHERE user_id = $1 AND room_id = $2
	`

	err := c.Db.QueryRowContext(ctx, query, userId, roomId).Scan(
		&userId,
		&roomId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to get chat room: %w", err)
	}

	return nil
}

func (c *ChatRepo) JoinChatRoom(ctx context.Context, req *entities.JoinChatRoomReq) error {
	query := `INSERT INTO users_rooms (room_id, user_id) VALUES ($1, $2)`
	_, err := c.Db.ExecContext(ctx, query, req.RoomId, req.UserId)
	if err != nil {
		return fmt.Errorf("failed to join chat room: %w", err)
	}

	return nil
}

func (c *ChatRepo) LeaveChatRoom(ctx context.Context, req *entities.JoinChatRoomReq) error {
	query := `DELETE FROM users_rooms WHERE room_id = $1 AND user_id = $2`
	_, err := c.Db.ExecContext(ctx, query, req.RoomId, req.UserId)
	if err != nil {
		return fmt.Errorf("failed to leave chat room: %w", err)
	}

	return nil
}

func (c *ChatRepo) SendMessage(ctx context.Context, req *entities.ChatMessage) error {
	query := `INSERT INTO messages (room_id, user_id, message) VALUES ($1, $2, $3)`
	_, err := c.Db.ExecContext(ctx, query, req.RoomId, req.Sender, req.Message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *ChatRepo) GetChatMessages(ctx context.Context, roomId int) ([]entities.ChatMessage, error) {
	var chatMessages []entities.ChatMessage
	query := `
	SELECT id, room_id, user_id, message, created_at FROM messages 
	WHERE room_id = $1 ORDER BY created_at ASC`
	err := c.Db.SelectContext(ctx, &chatMessages, query, roomId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entities.ChatMessage{}, nil
		}
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}

	return chatMessages, nil
}

func (c *ChatRepo) GetChatRoomUsers(ctx context.Context, roomId int) ([]entities.ChatUser, error) {
	var users []entities.ChatUser
	query := `SELECT users.id, users.username FROM users 
	JOIN users_rooms 
	ON users.id = users_rooms.user_id 
	WHERE users_rooms.room_id = $1`

	err := c.Db.SelectContext(ctx, &users, query, roomId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entities.ChatUser{}, nil
		}
		return nil, fmt.Errorf("failed to get chat room users: %w", err)
	}

	return users, nil
}
