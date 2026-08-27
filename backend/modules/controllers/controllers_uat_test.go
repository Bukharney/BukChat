package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukharney/bukchat/configs"
	"github.com/bukharney/bukchat/modules/controllers"
	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/gin-gonic/gin"
)

// Mock Auth Usecase
type mockAuthUsecase struct {
	loginFunc func(ctx context.Context, cfg *configs.Configs, req *entities.UsersCredentials) (*entities.UsersLoginRes, error)
}

func (m *mockAuthUsecase) Login(ctx context.Context, cfg *configs.Configs, req *entities.UsersCredentials) (*entities.UsersLoginRes, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, cfg, req)
	}
	return nil, nil
}

// Mock Users Usecase
type mockUsersUsecase struct {
	registerFunc       func(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error)
	changePasswordFunc func(ctx context.Context, req *entities.UsersChangePasswordReq) (*entities.UsersChangedRes, error)
	getUserDetailsFunc func(ctx context.Context, user entities.UsersClaims) (*entities.UsersDataRes, error)
	addFriendFunc      func(ctx context.Context, req *entities.FriendReq) (*entities.FriendRes, error)
	getFriendsFunc     func(ctx context.Context, userId int) ([]entities.FriendInfoRes, error)
	getFriendsReqFunc  func(ctx context.Context, userId int) ([]entities.FriendInfoRes, error)
}

func (m *mockUsersUsecase) Register(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUsersUsecase) ChangePassword(ctx context.Context, req *entities.UsersChangePasswordReq) (*entities.UsersChangedRes, error) {
	if m.changePasswordFunc != nil {
		return m.changePasswordFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUsersUsecase) GetUserDetails(ctx context.Context, user entities.UsersClaims) (*entities.UsersDataRes, error) {
	if m.getUserDetailsFunc != nil {
		return m.getUserDetailsFunc(ctx, user)
	}
	return nil, nil
}

func (m *mockUsersUsecase) DeleteAccount(ctx context.Context, user entities.UsersClaims) (*entities.UsersChangedRes, error) {
	return &entities.UsersChangedRes{Success: true}, nil
}

func (m *mockUsersUsecase) AddFriend(ctx context.Context, req *entities.FriendReq) (*entities.FriendRes, error) {
	if m.addFriendFunc != nil {
		return m.addFriendFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUsersUsecase) RejectFriend(ctx context.Context, userId int, friendUsername string) (*entities.UsersChangedRes, error) {
	return &entities.UsersChangedRes{Success: true}, nil
}

func (m *mockUsersUsecase) GetFriendsReq(ctx context.Context, userId int) ([]entities.FriendInfoRes, error) {
	if m.getFriendsReqFunc != nil {
		return m.getFriendsReqFunc(ctx, userId)
	}
	return nil, nil
}

func (m *mockUsersUsecase) GetFriends(ctx context.Context, userId int) ([]entities.FriendInfoRes, error) {
	if m.getFriendsFunc != nil {
		return m.getFriendsFunc(ctx, userId)
	}
	return nil, nil
}

// Mock Chat Usecase
type mockChatUsecase struct {
	createChatRoomFunc  func(ctx context.Context, req *entities.ChatRoom) error
	getChatRoomFunc     func(ctx context.Context, userId int, roomId int) error
	getChatMessagesFunc func(ctx context.Context, roomId int) ([]entities.ChatMessage, error)
}

func (m *mockChatUsecase) CreateChatRoom(ctx context.Context, req *entities.ChatRoom) error {
	if m.createChatRoomFunc != nil {
		return m.createChatRoomFunc(ctx, req)
	}
	return nil
}

func (m *mockChatUsecase) GetChatRoom(ctx context.Context, userId int, roomId int) error {
	if m.getChatRoomFunc != nil {
		return m.getChatRoomFunc(ctx, userId, roomId)
	}
	return nil
}

func (m *mockChatUsecase) GetChatMessages(ctx context.Context, roomId int) ([]entities.ChatMessage, error) {
	if m.getChatMessagesFunc != nil {
		return m.getChatMessagesFunc(ctx, roomId)
	}
	return nil, nil
}

func TestUAT_Auth_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &configs.Configs{}
	authUC := &mockAuthUsecase{
		loginFunc: func(ctx context.Context, cfg *configs.Configs, req *entities.UsersCredentials) (*entities.UsersLoginRes, error) {
			if req.Username == "testuser" && req.Password == "password123" {
				return &entities.UsersLoginRes{
					AccessToken: "mock-access-token",
				}, nil
			}
			return nil, apperrors.ErrInvalidCredentials
		},
	}

	controllers.NewAuthControllers(r.Group("/auth"), cfg, authUC)

	body, _ := json.Marshal(entities.UsersCredentials{
		Username: "testuser",
		Password: "password123",
	})

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got: %d", w.Code)
	}

	var res entities.UsersLoginRes
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	if res.AccessToken != "mock-access-token" {
		t.Errorf("Expected token 'mock-access-token', got: %s", res.AccessToken)
	}
}

func TestUAT_Auth_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &configs.Configs{}
	authUC := &mockAuthUsecase{
		loginFunc: func(ctx context.Context, cfg *configs.Configs, req *entities.UsersCredentials) (*entities.UsersLoginRes, error) {
			return nil, apperrors.ErrInvalidCredentials
		},
	}

	controllers.NewAuthControllers(r.Group("/auth"), cfg, authUC)

	body, _ := json.Marshal(entities.UsersCredentials{
		Username: "wronguser",
		Password: "wrongpassword",
	})

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected HTTP 401 Unauthorized, got: %d", w.Code)
	}
}

func TestUAT_User_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	authUC := &mockAuthUsecase{
		loginFunc: func(ctx context.Context, cfg *configs.Configs, req *entities.UsersCredentials) (*entities.UsersLoginRes, error) {
			return &entities.UsersLoginRes{AccessToken: "token-after-reg"}, nil
		},
	}
	userUC := &mockUsersUsecase{
		registerFunc: func(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
			return &entities.UsersRegisterRes{Id: 101, Username: req.Username}, nil
		},
	}

	controllers.NewUsersControllers(r.Group("/users"), userUC, authUC)

	regReq := entities.UsersRegisterReq{
		Username: "newuser",
		Password: "password123",
		Email:    "newuser@example.com",
	}
	body, _ := json.Marshal(regReq)

	req, _ := http.NewRequest("POST", "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got: %d", w.Code)
	}

	var res entities.UsersRegisterRes
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if res.Id != 101 || res.Username != "newuser" || res.AccessToken != "token-after-reg" {
		t.Errorf("Unexpected registration response: %+v", res)
	}
}

func TestUAT_User_Register_DuplicateConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	authUC := &mockAuthUsecase{}
	userUC := &mockUsersUsecase{
		registerFunc: func(ctx context.Context, req *entities.UsersRegisterReq) (*entities.UsersRegisterRes, error) {
			return nil, apperrors.ErrUsernameExists
		},
	}

	controllers.NewUsersControllers(r.Group("/users"), userUC, authUC)

	regReq := entities.UsersRegisterReq{
		Username: "existinguser",
		Password: "password123",
		Email:    "existing@example.com",
	}
	body, _ := json.Marshal(regReq)

	req, _ := http.NewRequest("POST", "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("Expected HTTP 409 Conflict, got: %d", w.Code)
	}
}

func TestUAT_Chat_CreateRoom_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	chatUC := &mockChatUsecase{
		createChatRoomFunc: func(ctx context.Context, req *entities.ChatRoom) error {
			if req.Name == "" {
				return errors.New("empty room name")
			}
			return nil
		},
	}

	controllers.NewChatControllers(r.Group("/chat"), nil, nil, chatUC)

	roomReq := entities.ChatRoom{
		Name: "Dev Channel",
	}
	body, _ := json.Marshal(roomReq)

	req, _ := http.NewRequest("POST", "/chat/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got: %d", w.Code)
	}
}
