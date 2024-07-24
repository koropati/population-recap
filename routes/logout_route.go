package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/controller"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/repository"
	"github.com/koropati/population-recap/usecase"
)

func NewLogoutRouter(cfg *SetupConfig, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(cfg.DB, domain.UserTable, cfg.Config.DefaultPageNumber, cfg.Config.DefaultPageSize)
	at := repository.NewAccessTokenRepository(cfg.DB, domain.AccessTokenTable, cfg.Config.DefaultPageNumber, cfg.Config.DefaultPageSize)
	rt := repository.NewRefreshTokenRepository(cfg.DB, domain.RefreshTokenTable, cfg.Config.DefaultPageNumber, cfg.Config.DefaultPageSize)
	lc := controller.LogoutController{
		UserUsecase:         usecase.NewUserUsecase(ur, cfg.Timeout),
		AccessTokenUsecase:  usecase.NewAccessTokenUsecase(at, cfg.Timeout),
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(rt, cfg.Timeout),
		Config:              cfg.Config,
		Cryptos:             cfg.Cryptos,
		Validator:           cfg.Validator,
	}

	group.GET("/logout", lc.Logout)
}
