package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/controller"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/repository"
	"github.com/koropati/population-recap/usecase"
)

func NewRegisterRouter(cfg *SetupConfig, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(cfg.DB, domain.UserTable, cfg.Config.DefaultPageNumber, cfg.Config.DefaultPageSize)
	sc := controller.RegisterController{
		UserUsecase: usecase.NewUserUsecase(ur, cfg.Timeout),
		Config:      cfg.Config,
		Cryptos:     cfg.Cryptos,
		Validator:   cfg.Validator,
	}

	group.GET("/register", sc.Index)
	group.POST("/register", sc.Register)
}
