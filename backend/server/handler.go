package server

import (
	"net/http"

	_usersController "github.com/bukharney/giga-chat/modules/controllers"
	_usersRepo "github.com/bukharney/giga-chat/modules/repositories"
	_usersUsecase "github.com/bukharney/giga-chat/modules/usecases"

	_authController "github.com/bukharney/giga-chat/modules/controllers"
	_authRepo "github.com/bukharney/giga-chat/modules/repositories"
	_authUsecase "github.com/bukharney/giga-chat/modules/usecases"

	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	usersGroup := v1.Group("/users")
	authGroup := v1.Group("/auth")

	usersRepo := _usersRepo.NewUsersRepo(s.DB)
	authRepo := _authRepo.NewAuthRepo(s.DB)

	authUsecase := _authUsecase.NewAuthUsecases(authRepo, usersRepo)
	usersUsecase := _usersUsecase.NewUsersUsecases(usersRepo)

	_usersController.NewUsersControllers(usersGroup, usersUsecase, authUsecase)
	_authController.NewAuthControllers(authGroup, s.Cfg, authUsecase)

	s.App.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Path Not Found"})
	})

	return nil
}
