package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/bootstrap"
	"github.com/koropati/population-recap/internal/cryptos"
	"github.com/koropati/population-recap/internal/validator"
)

type LandingPageController struct {
	Config    *bootstrap.Config
	Cryptos   cryptos.Cryptos
	Validator *validator.Validator
}

func (ctr *LandingPageController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "landing.tmpl", nil)
}
