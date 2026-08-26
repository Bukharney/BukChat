package usecases

import (
	"errors"
	"fmt"

	"github.com/bukharney/giga-chat/configs"
	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/bukharney/giga-chat/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecases struct {
	AuthRepo entities.AuthRepository
	UserRepo entities.UsersRepository
}

func NewAuthUsecases(authRepo entities.AuthRepository, userRepo entities.UsersRepository) entities.AuthUsecase {
	return &AuthUsecases{AuthRepo: authRepo, UserRepo: userRepo}
}

func (a *AuthUsecases) Login(cfg *configs.Configs, req *entities.UsersCredentials) (*entities.UsersLoginRes, error) {
	user, err := a.UserRepo.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := a.AuthRepo.SignUsersAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	res := &entities.UsersLoginRes{
		AccessToken: token,
	}
	return res, nil
}
