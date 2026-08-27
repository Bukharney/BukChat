package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/bukharney/bukchat/configs"
	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecases struct {
	AuthRepo entities.AuthRepository
	UserRepo entities.UsersRepository
}

func NewAuthUsecases(authRepo entities.AuthRepository, userRepo entities.UsersRepository) entities.AuthUsecase {
	return &AuthUsecases{AuthRepo: authRepo, UserRepo: userRepo}
}

func (a *AuthUsecases) Login(ctx context.Context, cfg *configs.Configs, req *entities.UsersCredentials) (*entities.UsersLoginRes, error) {
	user, err := a.UserRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := a.AuthRepo.SignUsersAccessToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	res := &entities.UsersLoginRes{
		AccessToken: token,
	}
	return res, nil
}
