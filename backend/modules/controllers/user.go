package controllers

import (
	"log/slog"
	"net/http"

	"github.com/bukharney/bukchat/configs"
	"github.com/bukharney/bukchat/middlewares"
	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/bukharney/bukchat/utils"
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

	ctx := c.Request.Context()
	res, err := u.UsersUsecase.Register(ctx, req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	token, err := u.AuthUsecase.Login(ctx, u.Cfg, user)
	if err != nil {
		slog.Error("Login after register failed", "error", err)
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

	res, err := u.UsersUsecase.ChangePassword(c.Request.Context(), req)
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

	res, err := u.UsersUsecase.GetUserDetails(c.Request.Context(), *user)
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

	res, err := u.UsersUsecase.DeleteAccount(c.Request.Context(), *user)
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

	res, err := u.UsersUsecase.AddFriend(c.Request.Context(), req)
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

	res, err := u.UsersUsecase.RejectFriend(c.Request.Context(), req.UserId, req.FriendUsername)
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

	res, err := u.UsersUsecase.GetFriendsReq(c.Request.Context(), user.Id)
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

	res, err := u.UsersUsecase.GetFriends(c.Request.Context(), user.Id)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
