package controllers

import (
	"net/http"

	"github.com/bukharney/bukchat/configs"
	"github.com/bukharney/bukchat/middlewares"
	"github.com/bukharney/bukchat/modules/entities"
	"github.com/bukharney/bukchat/pkg/apperrors"
	"github.com/bukharney/bukchat/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Cfg         *configs.Configs
	AuthUsecase entities.AuthUsecase
}

func NewAuthControllers(r gin.IRoutes, cfg *configs.Configs, authUsecase entities.AuthUsecase) {
	controllers := &AuthController{
		Cfg:         cfg,
		AuthUsecase: authUsecase,
	}

	r.POST("/login", controllers.Login)
	r.GET("/auth-test", middlewares.JwtAuthentication(), controllers.AuthTest)
	r.GET("/refresh-token", controllers.RefreshToken)
}

func (a *AuthController) Login(c *gin.Context) {
	req := new(entities.UsersCredentials)
	err := c.ShouldBind(req)
	if err != nil {
		utils.RespondWithError(c, apperrors.ErrBadRequest)
		return
	}

	res, err := a.AuthUsecase.Login(c.Request.Context(), a.Cfg, req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *AuthController) AuthTest(c *gin.Context) {
	tk, err := middlewares.GetUserByToken(c)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, tk)
}

func (a *AuthController) RefreshToken(c *gin.Context) {
	middlewares.RefreshToken(c)
}
