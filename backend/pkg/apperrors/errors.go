package apperrors

import (
	"errors"
	"net/http"
)

// Common Domain Errors
var (
	ErrNotFound             = errors.New("resource not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameExists       = errors.New("username already exists")
	ErrEmailExists          = errors.New("email already exists")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidCredentials   = errors.New("invalid username or password")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrMissingToken         = errors.New("missing authentication token")
	ErrFriendReqAlreadySent = errors.New("friend request already sent")
	ErrFriendAlreadyAdded   = errors.New("friend already added")
	ErrInternal             = errors.New("internal server error")
	ErrBadRequest           = errors.New("bad request")
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// MapErrorToStatus inspects an error and returns the corresponding HTTP status code and user-facing message.
func MapErrorToStatus(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code, appErr.Message
	}

	switch {
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound, ErrUserNotFound.Error()
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, ErrNotFound.Error()
	case errors.Is(err, ErrUsernameExists):
		return http.StatusConflict, ErrUsernameExists.Error()
	case errors.Is(err, ErrEmailExists):
		return http.StatusConflict, ErrEmailExists.Error()
	case errors.Is(err, ErrFriendReqAlreadySent):
		return http.StatusConflict, ErrFriendReqAlreadySent.Error()
	case errors.Is(err, ErrFriendAlreadyAdded):
		return http.StatusConflict, ErrFriendAlreadyAdded.Error()
	case errors.Is(err, ErrInvalidPassword):
		return http.StatusUnauthorized, ErrInvalidPassword.Error()
	case errors.Is(err, ErrInvalidCredentials):
		return http.StatusUnauthorized, ErrInvalidCredentials.Error()
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, ErrUnauthorized.Error()
	case errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized, ErrInvalidToken.Error()
	case errors.Is(err, ErrMissingToken):
		return http.StatusUnauthorized, ErrMissingToken.Error()
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, ErrBadRequest.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
