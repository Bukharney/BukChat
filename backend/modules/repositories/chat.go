package repositories

import (
	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/jmoiron/sqlx"
)

type ChatRepo struct {
	Db *sqlx.DB
}

func NewChatRepo(db *sqlx.DB) entities.ChatRepository {
	return &ChatRepo{Db: db}
}

func (c *ChatRepo) CreateChatRoom(req *entities.ChatRoom) error {
	query := `INSERT INTO chat_rooms (name) VALUES ($1)`
	_, err := c.Db.Exec(query, req.Name)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) GetChatRooms(userId int) ([]entities.ChatRoom, error) {
	var chatRooms []entities.ChatRoom
	query := `SELECT * FROM chat_rooms 
	JOIN chat_room_users 
	ON chat_rooms.id = chat_room_users.chat_room_id 
	WHERE chat_room_users.user_id = $1`

	err := c.Db.Select(&chatRooms, query, userId)
	if err != nil {
		return nil, err
	}

	return chatRooms, nil
}

func (c *ChatRepo) JoinChatRoom(req *entities.JoinChatRoomReq) error {
	query := `INSERT INTO chat_room_users (chat_room_id, user_id) VALUES ($1, $2)`
	_, err := c.Db.Exec(query, req.RoomId, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) LeaveChatRoom(req *entities.JoinChatRoomReq) error {
	query := `DELETE FROM chat_room_users WHERE chat_room_id = $1 AND user_id = $2`
	_, err := c.Db.Exec(query, req.RoomId, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) SendMessage(req *entities.ChatMessage) error {
	query := `INSERT INTO chat_messages (chat_room_id, sender, message) VALUES ($1, $2, $3)`
	_, err := c.Db.Exec(query, req.RoomId, req.Sender, req.Message)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) GetChatMessages(roomId int) ([]entities.ChatMessage, error) {
	var chatMessages []entities.ChatMessage
	query := `SELECT * FROM chat_messages WHERE chat_room_id = $1`
	err := c.Db.Select(&chatMessages, query, roomId)
	if err != nil {
		return nil, err
	}

	return chatMessages, nil
}

func (c *ChatRepo) GetChatRoomUsers(roomId int) ([]entities.ChatUser, error) {
	var users []entities.ChatUser
	query := `SELECT * FROM users 
	JOIN chat_room_users 
	ON users.id = chat_room_users.user_id 
	WHERE chat_room_users.chat_room_id = $1`

	err := c.Db.Select(&users, query, roomId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
