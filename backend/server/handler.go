package server

import (
	"net/http"

	_controller "github.com/bukharney/giga-chat/modules/controllers"
	_repo "github.com/bukharney/giga-chat/modules/repositories"
	_usecase "github.com/bukharney/giga-chat/modules/usecases"

	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	usersGroup := v1.Group("/users")
	authGroup := v1.Group("/auth")

	usersRepo := _repo.NewUsersRepo(s.DB)
	authRepo := _repo.NewAuthRepo(s.DB)

	authUsecase := _usecase.NewAuthUsecases(authRepo, usersRepo)
	usersUsecase := _usecase.NewUsersUsecases(usersRepo)

	_controller.NewUsersControllers(usersGroup, usersUsecase, authUsecase)
	_controller.NewAuthControllers(authGroup, s.Cfg, authUsecase)

	s.App.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Path Not Found"})
	})

	return nil
}
