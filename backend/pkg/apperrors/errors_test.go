package apperrors_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/bukharney/giga-chat/pkg/apperrors"
)

func TestMapErrorToStatus(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:         "Nil Error",
			err:          nil,
			expectedCode: http.StatusOK,
			expectedMsg:  "",
		},
		{
			name:         "ErrUserNotFound Sentinel",
			err:          apperrors.ErrUserNotFound,
			expectedCode: http.StatusNotFound,
			expectedMsg:  "user not found",
		},
		{
			name:         "Wrapped ErrUserNotFound",
			err:          fmt.Errorf("repo get failed: %w", apperrors.ErrUserNotFound),
			expectedCode: http.StatusNotFound,
			expectedMsg:  "user not found",
		},
		{
			name:         "ErrUsernameExists Conflict",
			err:          apperrors.ErrUsernameExists,
			expectedCode: http.StatusConflict,
			expectedMsg:  "username already exists",
		},
		{
			name:         "ErrInvalidCredentials Unauthorized",
			err:          apperrors.ErrInvalidCredentials,
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "invalid username or password",
		},
		{
			name:         "Custom AppError",
			err:          apperrors.NewAppError(http.StatusForbidden, "access forbidden", errors.New("underlying detail")),
			expectedCode: http.StatusForbidden,
			expectedMsg:  "access forbidden",
		},
		{
			name:         "Unknown Internal Error",
			err:          errors.New("something went wrong DB connection failed"),
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, msg := apperrors.MapErrorToStatus(tt.err)
			if code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, code)
			}
			if msg != tt.expectedMsg {
				t.Errorf("expected msg %q, got %q", tt.expectedMsg, msg)
			}
		})
	}
}
