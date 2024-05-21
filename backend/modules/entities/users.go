package entities

import (
	"github.com/golang-jwt/jwt/v4"
)

type UsersUsecase interface {
	Register(req *UsersRegisterReq) (*UsersRegisterRes, error)
	ChangePassword(req *UsersChangePasswordReq) (*UsersChangedRes, error)
	GetUserDetails(user UsersClaims) (*UsersDataRes, error)
	DeleteAccount(user UsersClaims) (*UsersChangedRes, error)
	AddFriend(req *FriendReq) (*FriendRes, error)
	GetFriendsReq(userId int) ([]FriendInfoRes, error)
	GetFriends(userId int) ([]FriendInfoRes, error)
}

type UsersRepository interface {
	Register(req *UsersRegisterReq) (*UsersRegisterRes, error)
	GetUserByUsername(username string) (*UsersPassport, error)
	ChangePassword(req *UsersChangePasswordReq) (*UsersChangedRes, error)
	DeleteAccount(user_id int) (*UsersChangedRes, error)
	AddFriend(req *FriendReq) (*FriendRes, error)
	GetFriendsReq(user_id int) ([]FriendInfoRes, error)
	GetFriendReq(user_id int, friend_id int) (*FriendRes, error)
	GetFriends(user_id int) ([]FriendInfoRes, error)
	AcceptFriendReq(user_id int, friend_id int, room_id int) (*FriendRes, error)
}

type UsersCredentials struct {
	Username string `json:"username" db:"username" form:"username" binding:"required"`
	Password string `json:"password" db:"password" form:"password" binding:"required"`
}

type UsersPassport struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
}

type UsersDataRes struct {
	Id          int    `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	Email       string `json:"email" db:"email"`
	AccessToken string `json:"token"`
}

type UsersClaims struct {
	Id       int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UsersRegisterReq struct {
	Username string `json:"username" db:"username" binding:"required"`
	Password string `json:"password" db:"password" binding:"required"`
	Email    string `json:"email" db:"email" binding:"required"`
}

type UsersChangePasswordReq struct {
	Id          int    `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	OldPassword string `json:"old_password" db:"old_password" binding:"required"`
	NewPassword string `json:"new_password" db:"new_password" binding:"required"`
}

type UsersRegisterRes struct {
	Id          int    `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	AccessToken string `json:"token"`
}

type UsersLoginRes struct {
	AccessToken string `json:"token"`
}

type UsersChangedRes struct {
	Success bool `json:"success"`
}

type FriendReq struct {
	UserId         int    `json:"user_id"`
	FriendId       int    `json:"friend_id"`
	FriendUsername string `json:"username" binding:"required"`
	Status         int    `json:"status"`
}

type FriendRes struct {
	UserId   int    `json:"user_id" db:"from_user_id"`
	FriendId int    `json:"friend_id" db:"to_user_id"`
	Status   int    `json:"status" db:"status"`
	Created  string `json:"created_at" db:"created_at"`
}

type FriendInfoRes struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Status   int    `json:"status"`
	RoomId   int    `json:"room_id" db:"room_id"`
}
