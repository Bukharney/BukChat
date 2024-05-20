package entities

type ChatRepository interface {
	CreateChatRoom(req *ChatRoom) error
}

type ChatUsecase interface {
	CreateChatRoom(req *ChatRoom) error
}

type ChatRoom struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ChatMessage struct {
	Id      int    `json:"id"`
	RoomId  int    `json:"RoomId"`
	Sender  int    `json:"sender"`
	Message string `json:"message"`
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
