package utils

import (
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RespondWithError(c *gin.Context, err error) {
	statusCode, message := apperrors.MapErrorToStatus(err)
	c.JSON(statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func RespondWithCustomError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}
