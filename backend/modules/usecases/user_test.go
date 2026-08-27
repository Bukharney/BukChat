package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
)

type mockUserRepo struct {
	registerFunc          func(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error)
	getUserByUsernameFunc func(ctx context.Context, username string) (*entities.UsersPassport, error)
}

func (m *mockUserRepo) Register(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUserRepo) GetUserByUsername(ctx context.Context, username string) (*entities.UsersPassport, error) {
	if m.getUserByUsernameFunc != nil {
		return m.getUserByUsernameFunc(ctx, username)
	}
	return nil, nil
}

func (m *mockUserRepo) ChangePassword(ctx context.Context, req *entities.UsersChangePasswordReq) (*entities.UsersChangedRes, error) {
	return nil, nil
}

func (m *mockUserRepo) DeleteAccount(ctx context.Context, user_id int) (*entities.UsersChangedRes, error) {
	return nil, nil
}

func (m *mockUserRepo) AddFriend(ctx context.Context, req *entities.FriendReq) (*entities.FriendRes, error) {
	return nil, nil
}

func (m *mockUserRepo) GetFriendsReq(ctx context.Context, user_id int) ([]entities.FriendInfoRes, error) {
	return nil, nil
}

func (m *mockUserRepo) GetFriendReq(ctx context.Context, user_id int, friend_id int) (*entities.FriendRes, error) {
	return nil, nil
}

func (m *mockUserRepo) GetFriends(ctx context.Context, user_id int) ([]entities.FriendInfoRes, error) {
	return nil, nil
}

func (m *mockUserRepo) AcceptFriendReq(ctx context.Context, user_id int, friend_id int, room_id int) (*entities.FriendRes, error) {
	return nil, nil
}

func (m *mockUserRepo) RejectFriend(ctx context.Context, user_id int, friend_id int) (*entities.UsersChangedRes, error) {
	return nil, nil
}

func TestUsersUsecases_Register_Success(t *testing.T) {
	repo := &mockUserRepo{
		registerFunc: func(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
			return &entities.UsersRegisterRes{
				Id:       1,
				Username: req.Username,
			}, nil
		},
	}

	uc := NewUsersUsecases(repo, nil, nil)
	ctx := context.Background()

	res, err := uc.Register(ctx, &entities.UsersRegisterReq{
		Username: "alice",
		Password: "password123",
		Email:    "alice@example.com",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if res.Id != 1 || res.Username != "alice" {
		t.Errorf("Unexpected result: %+v", res)
	}
}

func TestUsersUsecases_Register_DuplicateUser(t *testing.T) {
	repo := &mockUserRepo{
		registerFunc: func(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
			return nil, apperrors.ErrUsernameExists
		},
	}

	uc := NewUsersUsecases(repo, nil, nil)
	ctx := context.Background()

	_, err := uc.Register(ctx, &entities.UsersRegisterReq{
		Username: "existinguser",
		Password: "password123",
		Email:    "existing@example.com",
	})

	if !errors.Is(err, apperrors.ErrUsernameExists) {
		t.Fatalf("Expected ErrUsernameExists, got: %v", err)
	}
}
