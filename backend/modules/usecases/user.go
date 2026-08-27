package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/bukharney/bukchat/server/ws"
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

func (a *UsersUsecases) Register(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword

	user, err := a.UsersRepo.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UsersUsecases) ChangePassword(ctx context.Context, req *entities.UsersChangePasswordReq) (*entities.UsersChangedRes, error) {
	user, err := a.UsersRepo.GetUserByUsername(ctx, req.Username)
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

	res, err := a.UsersRepo.ChangePassword(ctx, req)
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

func (a *UsersUsecases) GetUserByUsername(ctx context.Context, username string) (*entities.UsersPassport, error) {
	user, err := a.UsersRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UsersUsecases) GetUserDetails(ctx context.Context, user entities.UsersClaims) (*entities.UsersDataRes, error) {
	res, err := a.UsersRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	return &entities.UsersDataRes{
		Id:       res.Id,
		Username: res.Username,
		Email:    res.Email,
	}, nil
}

func (a *UsersUsecases) DeleteAccount(ctx context.Context, user entities.UsersClaims) (*entities.UsersChangedRes, error) {
	if user.Id == 0 {
		return nil, apperrors.ErrUserNotFound
	}

	res, err := a.UsersRepo.DeleteAccount(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *UsersUsecases) AddFriend(ctx context.Context, req *entities.FriendReq) (*entities.FriendRes, error) {
	friendId, err := a.GetUserByUsername(ctx, req.FriendUsername)
	if err != nil {
		return nil, err
	}

	req.FriendId = friendId.Id
	status, err := a.UsersRepo.GetFriendReq(ctx, req.UserId, friendId.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			req.Status = 0
			res, err := a.UsersRepo.AddFriend(ctx, req)
			if err != nil {
				return nil, err
			}
			if a.Hub != nil {
				a.Hub.SendToUser(friendId.Id, "friend_request", map[string]interface{}{
					"from_user_id":  req.UserId,
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
		id, err := a.ChatRepo.CreateChatRoom(ctx, &entities.ChatRoom{
			Name: friendId.Username,
		})
		if err != nil {
			return nil, err
		}

		res, err := a.UsersRepo.AcceptFriendReq(ctx, status.UserId, status.FriendId, id)
		if err != nil {
			return nil, err
		}

		err = a.ChatRepo.JoinChatRoom(ctx, &entities.JoinChatRoomReq{
			UserId: friendId.Id,
			RoomId: id,
		})
		if err != nil {
			return nil, err
		}

		err = a.ChatRepo.JoinChatRoom(ctx, &entities.JoinChatRoomReq{
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

func (a *UsersUsecases) RejectFriend(ctx context.Context, userId int, FriendUsername string) (*entities.UsersChangedRes, error) {
	friend, err := a.GetUserByUsername(ctx, FriendUsername)
	if err != nil {
		return nil, err
	}
	res, err := a.UsersRepo.RejectFriend(ctx, userId, friend.Id)
	if err != nil {
		return nil, err
	}

	if a.Hub != nil {
		a.Hub.SendToUser(friend.Id, "friend_rejected", res)
	}

	return res, nil
}

func (a *UsersUsecases) GetFriendsReq(ctx context.Context, userId int) ([]entities.FriendInfoRes, error) {
	res, err := a.UsersRepo.GetFriendsReq(ctx, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *UsersUsecases) GetFriends(ctx context.Context, userId int) ([]entities.FriendInfoRes, error) {
	res, err := a.UsersRepo.GetFriends(ctx, userId)
	if err != nil {
		return nil, err
	}

	for i := range res {
		if a.Hub != nil {
			res[i].IsOnline = a.Hub.IsUserOnline(res[i].Id)
		}
	}

	return res, nil
}
