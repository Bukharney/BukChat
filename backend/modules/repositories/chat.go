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

func (c *ChatRepo) CreateChatRoom(req *entities.ChatRoom) (int, error) {
	query := `INSERT INTO rooms (name) VALUES ($1) RETURNING id`
	var roomId int
	err := c.Db.QueryRow(query, req.Name).Scan(&roomId)
	if err != nil {
		return 0, err
	}

	return roomId, nil
}

func (c *ChatRepo) GetChatRoom(userId int, roomId int) error {
	query := `SELECT * FROM users_rooms
	WHERE user_id = $1 AND room_id = $2
	`

	err := c.Db.QueryRow(query, userId, roomId).Scan(
		&userId,
		&roomId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) JoinChatRoom(req *entities.JoinChatRoomReq) error {
	query := `INSERT INTO users_rooms (room_id, user_id) VALUES ($1, $2)`
	_, err := c.Db.Exec(query, req.RoomId, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) LeaveChatRoom(req *entities.JoinChatRoomReq) error {
	query := `DELETE FROM users_rooms WHERE room_id = $1 AND user_id = $2`
	_, err := c.Db.Exec(query, req.RoomId, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) SendMessage(req *entities.ChatMessage) error {
	query := `INSERT INTO messages (room_id, user_id, message) VALUES ($1, $2, $3)`
	_, err := c.Db.Exec(query, req.RoomId, req.Sender, req.Message)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepo) GetChatMessages(roomId int) ([]entities.ChatMessage, error) {
	var chatMessages []entities.ChatMessage
	query := `
	SELECT * FROM messages 
	WHERE room_id = $1`
	err := c.Db.Select(&chatMessages, query, roomId)
	if err != nil {
		return nil, err
	}

	return chatMessages, nil
}

func (c *ChatRepo) GetChatRoomUsers(roomId int) ([]entities.ChatUser, error) {
	var users []entities.ChatUser
	query := `SELECT * FROM users 
	JOIN users_rooms 
	ON users.id = users_rooms.user_id 
	WHERE users_rooms.room_id = $1`

	err := c.Db.Select(&users, query, roomId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
