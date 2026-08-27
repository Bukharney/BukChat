package controllers

import (
	"net/http"
	"strconv"

	"github.com/bukharney/bukchat/configs"
	"github.com/bukharney/bukchat/middlewares"
	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/bukharney/bukchat/utils"
	"github.com/gin-gonic/gin"
)

type ChatController struct {
	Cfg          *configs.Configs
	UsersUsecase entities.UsersUsecase
	AuthUsecase  entities.AuthUsecase
	ChatUsecase  entities.ChatUsecase
}

func NewChatControllers(r gin.IRoutes, usersUsecase entities.UsersUsecase, authUsecase entities.AuthUsecase, chatUsecase entities.ChatUsecase) {
	controllers := &ChatController{
		UsersUsecase: usersUsecase,
		AuthUsecase:  authUsecase,
		ChatUsecase:  chatUsecase,
	}

	r.POST("/", controllers.CreateChatRoom)
	r.GET("/:roomId", middlewares.JwtAuthentication(), controllers.GetChatMessages)
}

func (c *ChatController) CreateChatRoom(ctx *gin.Context) {
	var req entities.ChatRoom
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondWithError(ctx, apperrors.ErrBadRequest)
		return
	}

	err = c.ChatUsecase.CreateChatRoom(ctx.Request.Context(), &req)
	if err != nil {
		utils.RespondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "chat room created successfully"})
}

func (c *ChatController) GetChatMessages(ctx *gin.Context) {
	user, err := middlewares.GetUserByToken(ctx)
	if err != nil {
		utils.RespondWithError(ctx, err)
		return
	}
	roomId := ctx.Param("roomId")
	rid, err := strconv.Atoi(roomId)
	if err != nil {
		utils.RespondWithError(ctx, apperrors.ErrBadRequest)
		return
	}

	reqCtx := ctx.Request.Context()
	err = c.ChatUsecase.GetChatRoom(reqCtx, user.Id, rid)
	if err != nil {
		utils.RespondWithError(ctx, err)
		return
	}

	messages, err := c.ChatUsecase.GetChatMessages(reqCtx, rid)
	if err != nil {
		utils.RespondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
