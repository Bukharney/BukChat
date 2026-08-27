package entities

import (
	"context"

	"github.com/bukharney/bukchat/configs"
)

type AuthRepository interface {
	SignUsersAccessToken(ctx context.Context, req *UsersPassport) (string, error)
}

type AuthUsecase interface {
	Login(ctx context.Context, cfg *configs.Configs, req *UsersCredentials) (*UsersLoginRes, error)
}
