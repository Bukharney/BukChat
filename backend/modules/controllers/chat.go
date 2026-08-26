package controllers

import (
	"net/http"
	"strconv"

	"github.com/bukharney/giga-chat/configs"
	"github.com/bukharney/giga-chat/middlewares"
	"github.com/bukharney/giga-chat/modules/entities"
	"github.com/bukharney/giga-chat/pkg/apperrors"
	"github.com/bukharney/giga-chat/utils"
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

	err = c.ChatUsecase.CreateChatRoom(&req)
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

	err = c.ChatUsecase.GetChatRoom(user.Id, rid)
	if err != nil {
		utils.RespondWithError(ctx, err)
		return
	}

	messages, err := c.ChatUsecase.GetChatMessages(rid)
	if err != nil {
		utils.RespondWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
