package usecases

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/bukharney/giga-chat/pkg/apperrors"
	"github.com/bukharney/giga-chat/server/ws"
	"golang.org/x/crypto/bcrypt"
)

type UsersUsecases struct {
	UsersRepo entities.UsersRepository
	ChatRepo  entities.ChatRepository
	Hub       *ws.Hub
}

func NewUsersUsecases(usersRepo entities.UsersRepository, chatRepo entities.ChatRepository, hub *ws.Hub) entities.UsersUsecase {
	return &UsersUsecases{UsersRepo: usersRepo, ChatRepo: chatRepo, Hub: hub}
}

func (a *UsersUsecases) Register(req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword

	user, err := a.UsersRepo.Register(req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UsersUsecases) ChangePassword(req *entities.UsersChangePasswordReq) (*entities.UsersChangedRes, error) {
	user, err := a.UsersRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return nil, apperrors.ErrInvalidPassword
	}

	req.NewPassword, err = hashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}

	res, err := a.UsersRepo.ChangePassword(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func (a *UsersUsecases) GetUserByUsername(username string) (*entities.UsersPassport, error) {
	user, err := a.UsersRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UsersUsecases) GetUserDetails(user entities.UsersClaims) (*entities.UsersDataRes, error) {
	res, err := a.UsersRepo.GetUserByUsername(user.Username)
	if err != nil {
		return nil, err
	}

	return &entities.UsersDataRes{
		Id:       res.Id,
		Username: res.Username,
		Email:    res.Email,
	}, nil
}

func (a *UsersUsecases) DeleteAccount(user entities.UsersClaims) (*entities.UsersChangedRes, error) {
	if user.Id == 0 {
		return nil, apperrors.ErrUserNotFound
	}

	res, err := a.UsersRepo.DeleteAccount(user.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *UsersUsecases) AddFriend(req *entities.FriendReq) (*entities.FriendRes, error) {
	friendId, err := a.GetUserByUsername(req.FriendUsername)
	if err != nil {
		return nil, err
	}

	req.FriendId = friendId.Id
	status, err := a.UsersRepo.GetFriendReq(req.UserId, friendId.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			req.Status = 0
			res, err := a.UsersRepo.AddFriend(req)
			if err != nil {
				return nil, err
			}
			if a.Hub != nil {
				a.Hub.SendToUser(friendId.Id, "friend_request", map[string]interface{}{
					"from_user_id": req.UserId,
					"from_username": req.FriendUsername,
				})
			}
			return res, nil
		}
		return nil, err
	}

	if status.Status == 0 && status.UserId == req.UserId && status.FriendId == friendId.Id {
		return nil, apperrors.ErrFriendReqAlreadySent
	} else if status.Status == 1 && status.UserId == req.UserId && status.FriendId == friendId.Id {
		return nil, apperrors.ErrFriendAlreadyAdded
	} else if status.Status == 0 && status.UserId == friendId.Id && status.FriendId == req.UserId {
		id, err := a.ChatRepo.CreateChatRoom(&entities.ChatRoom{
			Name: friendId.Username,
		})
		if err != nil {
			return nil, err
		}

		res, err := a.UsersRepo.AcceptFriendReq(status.UserId, status.FriendId, id)
		if err != nil {
			return nil, err
		}

		err = a.ChatRepo.JoinChatRoom(&entities.JoinChatRoomReq{
			UserId: friendId.Id,
			RoomId: id,
		})
		if err != nil {
			return nil, err
		}

		err = a.ChatRepo.JoinChatRoom(&entities.JoinChatRoomReq{
			UserId: req.UserId,
			RoomId: id,
		})
		if err != nil {
			return nil, err
		}

		if a.Hub != nil {
			a.Hub.SendToUser(status.UserId, "friend_accepted", res)
			a.Hub.SendToUser(status.FriendId, "friend_accepted", res)
		}

		return res, nil
	} else if status.Status == 1 && status.UserId == friendId.Id && status.FriendId == req.UserId {
		return nil, apperrors.ErrFriendAlreadyAdded
	} else {
		return nil, apperrors.ErrInternal
	}
}

func (a *UsersUsecases) RejectFriend(userId int, FriendUsername string) (*entities.UsersChangedRes, error) {
	friend, err := a.GetUserByUsername(FriendUsername)
	if err != nil {
		return nil, err
	}
	res, err := a.UsersRepo.RejectFriend(userId, friend.Id)
	if err != nil {
		return nil, err
	}

	if a.Hub != nil {
		a.Hub.SendToUser(friend.Id, "friend_rejected", res)
	}

	return res, nil
}

func (a *UsersUsecases) GetFriendsReq(userId int) ([]entities.FriendInfoRes, error) {
	res, err := a.UsersRepo.GetFriendsReq(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *UsersUsecases) GetFriends(userId int) ([]entities.FriendInfoRes, error) {
	res, err := a.UsersRepo.GetFriends(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
