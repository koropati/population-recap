package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/controller"
)

func NewLandingPageRouter(cfg *SetupConfig, group *gin.RouterGroup) {
	sc := controller.LandingPageController{
		Config:    cfg.Config,
		Cryptos:   cfg.Cryptos,
		Validator: cfg.Validator,
	}

	group.GET("/", sc.Index)
}
