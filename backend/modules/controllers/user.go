package controllers

import (
	"log"
	"net/http"

	"github.com/bukharney/giga-chat/configs"
	"github.com/bukharney/giga-chat/middlewares"
	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/bukharney/giga-chat/pkg/apperrors"
	"github.com/bukharney/giga-chat/utils"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	Cfg          *configs.Configs
	UsersUsecase entities.UsersUsecase
	AuthUsecase  entities.AuthUsecase
}

func NewUsersControllers(r gin.IRoutes, usersUsecase entities.UsersUsecase, authUsecase entities.AuthUsecase) {
	controllers := &UsersController{
		UsersUsecase: usersUsecase,
		AuthUsecase:  authUsecase,
	}

	r.GET("/", middlewares.JwtAuthentication(), controllers.GetUserDetails)
	r.GET("/friends-request", middlewares.JwtAuthentication(), controllers.GetFriendsReq)
	r.GET("/friends", middlewares.JwtAuthentication(), controllers.GetFriends)
	r.POST("/", controllers.Register)
	r.POST("/add-friend", middlewares.JwtAuthentication(), controllers.AddFriend)
	r.POST("/reject-friend", middlewares.JwtAuthentication(), controllers.RejectFriend)
	r.DELETE("/", middlewares.JwtAuthentication(), controllers.DeleteAccount)
	r.PATCH("/", middlewares.JwtAuthentication(), controllers.ChangePassword)
}

func (u *UsersController) Register(c *gin.Context) {
	req := new(entities.UsersRegisterReq)
	err := c.ShouldBind(req)
	if err != nil {
		utils.RespondWithError(c, apperrors.ErrBadRequest)
		return
	}

	user := &entities.UsersCredentials{
		Username: req.Username,
		Password: req.Password,
	}

	res, err := u.UsersUsecase.Register(req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	token, err := u.AuthUsecase.Login(u.Cfg, user)
	if err != nil {
		log.Println("login after register failed:", err)
		utils.RespondWithError(c, err)
		return
	}

	res.AccessToken = token.AccessToken

	c.JSON(http.StatusOK, res)
}

func (u *UsersController) ChangePassword(c *gin.Context) {
	claims, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	req := new(entities.UsersChangePasswordReq)
	err = c.ShouldBind(req)
	if err != nil {
		utils.RespondWithError(c, apperrors.ErrBadRequest)
		return
	}
	req.Username = claims.Username
	req.Id = claims.Id

	res, err := u.UsersUsecase.ChangePassword(req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *UsersController) GetUserDetails(c *gin.Context) {
	user, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	res, err := u.UsersUsecase.GetUserDetails(*user)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *UsersController) DeleteAccount(c *gin.Context) {
	user, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	res, err := u.UsersUsecase.DeleteAccount(*user)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *UsersController) AddFriend(c *gin.Context) {
	user, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	req := new(entities.FriendReq)
	err = c.ShouldBind(req)
	if err != nil {
		utils.RespondWithError(c, apperrors.ErrBadRequest)
		return
	}

	req.UserId = user.Id

	res, err := u.UsersUsecase.AddFriend(req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *UsersController) RejectFriend(c *gin.Context) {
	user, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	req := new(entities.FriendReq)
	err = c.ShouldBind(req)
	if err != nil {
		utils.RespondWithError(c, apperrors.ErrBadRequest)
		return
	}

	req.UserId = user.Id

	res, err := u.UsersUsecase.RejectFriend(req.UserId, req.FriendUsername)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *UsersController) GetFriendsReq(c *gin.Context) {
	user, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	res, err := u.UsersUsecase.GetFriendsReq(user.Id)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *UsersController) GetFriends(c *gin.Context) {
	user, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	res, err := u.UsersUsecase.GetFriends(user.Id)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
