package routes

import (
	"time"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/bootstrap"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/internal/cryptos"
	"github.com/koropati/population-recap/internal/mailer"
	"github.com/koropati/population-recap/internal/validator"
	"github.com/koropati/population-recap/middleware"
	"github.com/koropati/population-recap/repository"
	"github.com/koropati/population-recap/usecase"
	"gorm.io/gorm"
)

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleStaff      = "staff"
	DeviceSwitch   = "switch"
	DeviceHps      = "hps"
	DeviceIr       = "ir"
	DeviceAc       = "ac"
)

type SetupConfig struct {
	Config         *bootstrap.Config
	Timeout        time.Duration
	DB             *gorm.DB
	CasbinEnforcer *casbin.Enforcer
	Cryptos        cryptos.Cryptos
	Gin            *gin.Engine
	Validator      *validator.Validator
	Mailer         mailer.Mailer
}

func Setup(config *SetupConfig) {

	config.Gin.Static("assets", "./templates/assets")
	config.Gin.LoadHTMLGlob("./templates/*.tmpl")

	at := repository.NewAccessTokenRepository(config.DB, domain.AccessTokenTable, config.Config.DefaultPageNumber, config.Config.DefaultPageSize)
	rt := repository.NewRefreshTokenRepository(config.DB, domain.RefreshTokenTable, config.Config.DefaultPageNumber, config.Config.DefaultPageSize)

	// All Public APIs
	publicRouter := config.Gin.Group("/")
	publicRouter.Use(middleware.AuthPublicMiddleware(config.Config.AccessTokenSecret, config.CasbinEnforcer, config.Cryptos, usecase.NewAccessTokenUsecase(at, config.Timeout), usecase.NewRefreshTokenUsecase(rt, config.Timeout)))
	NewLandingPageRouter(config, publicRouter)
	NewRegisterRouter(config, publicRouter)
	NewLoginRouter(config, publicRouter)
	NewLogoutRouter(config, publicRouter)
	NewForgotPasswordRouter(config, publicRouter)

	privateRouter := config.Gin.Group("/")
	privateRouter.Use(middleware.AuthMiddleware(config.Config.AccessTokenSecret, config.CasbinEnforcer, config.Cryptos, usecase.NewAccessTokenUsecase(at, config.Timeout), usecase.NewRefreshTokenUsecase(rt, config.Timeout)))
	NewDashboardPageRouter(config, privateRouter)

}
