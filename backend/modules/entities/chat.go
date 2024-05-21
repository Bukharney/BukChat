package entities

type ChatRepository interface {
	CreateChatRoom(req *ChatRoom) (int, error)
	GetChatRoom(userId int, roomId int) error
	JoinChatRoom(req *JoinChatRoomReq) error
	SendMessage(req *ChatMessage) error
	GetChatMessages(roomId int) ([]ChatMessage, error)
}

type ChatUsecase interface {
	CreateChatRoom(req *ChatRoom) error
	GetChatMessages(roomId int) ([]ChatMessage, error)
	GetChatRoom(userId int, roomId int) error
}

type ChatRoom struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ChatMessage struct {
	Id        int    `json:"id"`
	RoomId    int    `json:"room_id" db:"room_id"`
	Sender    int    `json:"user_id" db:"user_id"`
	Message   string `json:"message" db:"message"`
	Timestamp string `json:"timestamp" db:"created_at"`
}

type JoinChatRoomReq struct {
	UserId int `json:"user_id"`
	RoomId int `json:"room_id"`
}

type ChatRoomUsersRes struct {
	Chatrooms []ChatRoom `json:"chatrooms"`
}

type ChatUsersRes struct {
	Users []ChatUser `json:"users"`
}

type ChatUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
