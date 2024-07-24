package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/controller"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/repository"
	"github.com/koropati/population-recap/usecase"
)

func NewForgotPasswordRouter(cfg *SetupConfig, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(cfg.DB, domain.UserTable, cfg.Config.DefaultPageNumber, cfg.Config.DefaultPageSize)
	fpt := repository.NewForgotPasswordTokenRepository(cfg.DB, domain.ForgotPasswordTokenTable, cfg.Config.DefaultPageNumber, cfg.Config.DefaultPageSize)

	lc := controller.ForgotPasswordController{
		UserUsecase:                usecase.NewUserUsecase(ur, cfg.Timeout),
		ForgotPasswordTokenUsecase: usecase.NewForgotPasswordTokenUsecase(fpt, cfg.Timeout),
		Config:                     cfg.Config,
		Cryptos:                    cfg.Cryptos,
		Validator:                  cfg.Validator,
		Mailer:                     cfg.Mailer,
	}

	group.GET("/forgot-password", lc.Index)
	group.POST("/forgot-password", lc.ForgotPassword)
	group.GET("/forgot-password/verify", lc.VerifyForgotPassword)
	group.GET("/reset-password", lc.ResetPassword)
	group.POST("/reset-password", lc.ResetPasswordConfirm)
	// /verify-forgot-password
}
