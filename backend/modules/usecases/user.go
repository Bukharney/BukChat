package usecases

import (
	"errors"

	"github.com/bukharney/giga-chat/modules/entities"
	"golang.org/x/crypto/bcrypt"
)

type UsersUsecases struct {
	UsersRepo entities.UsersRepository
}

func NewUsersUsecases(usersRepo entities.UsersRepository) entities.UsersUsecase {
	return &UsersUsecases{UsersRepo: usersRepo}
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
		return nil, errors.New("error, user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return nil, errors.New("error, password is invalid")
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
		return "", err
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
		return nil, errors.New("error, user not found")
	}

	res, err := a.UsersRepo.DeleteAccount(user.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *UsersUsecases) AddFriend(req *entities.FriendReq) (*entities.FriendRes, error) {
	status, err := a.UsersRepo.GetFriendReq(req.UserId, req.FriendId)
	if err != nil {
		return nil, err
	}
	if status == 0 {
		req.Status = 1
		res, err := a.UsersRepo.AddFriend(req)
		if err != nil {
			return nil, err
		}

		return res, nil
	} else {
		req.Status = 0
		res, err := a.UsersRepo.AddFriend(req)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}
