package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/bukharney/bukchat/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		notAuth := []string{"/login"}
		requestPath := c.Request.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				c.Next()
				return
			}
		}

		tokenHeader := c.Request.Header.Get("Authorization")

		if tokenHeader == "" {
			utils.RespondWithError(c, apperrors.ErrMissingToken)
			c.Abort()
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			utils.RespondWithError(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenPart := splitted[1]
		tk := &entities.UsersClaims{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			utils.RespondWithError(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		if !token.Valid {
			utils.RespondWithError(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserByToken(c *gin.Context) (*entities.UsersClaims, error) {
	tokenHeader := c.Request.Header.Get("Authorization")

	if tokenHeader == "" {
		return nil, apperrors.ErrMissingToken
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		return nil, apperrors.ErrInvalidToken
	}

	tokenPart := splitted[1]

	tk := &entities.UsersClaims{}

	if tokenPart == "" {
		return nil, apperrors.ErrMissingToken
	}

	_, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	return tk, nil
}

func GetUserToken(tokenPart string) (*entities.UsersClaims, error) {
	tk := &entities.UsersClaims{}

	if tokenPart == "" {
		return nil, apperrors.ErrMissingToken
	}

	_, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	return tk, nil
}

func RefreshToken(c *gin.Context) {
	tokenHeader := c.Request.Header.Get("Authorization")
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		utils.RespondWithError(c, apperrors.ErrInvalidToken)
		return
	}
	tokenPart := splitted[1]
	tk := &entities.UsersClaims{}
	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		utils.RespondWithError(c, apperrors.ErrInvalidToken)
		return
	}

	if token.Valid {
		claims := token.Claims.(*entities.UsersClaims)
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
		claims.NotBefore = jwt.NewNumericDate(time.Now())
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			utils.RespondWithError(c, apperrors.ErrInternal)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": ss,
		})

	} else {
		utils.RespondWithError(c, apperrors.ErrInvalidToken)
	}
}
