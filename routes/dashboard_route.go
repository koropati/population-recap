package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/controller"
)

func NewDashboardPageRouter(cfg *SetupConfig, group *gin.RouterGroup) {
	sc := controller.DashboardController{
		Config:    cfg.Config,
		Cryptos:   cfg.Cryptos,
		Validator: cfg.Validator,
	}

	group.GET("/dashboard", sc.Index)
}
