package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/bootstrap"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/internal/cryptos"
	"github.com/koropati/population-recap/internal/validator"
)

type RegisterController struct {
	UserUsecase domain.UserUsecase
	Config      *bootstrap.Config
	Cryptos     cryptos.Cryptos
	Validator   *validator.Validator
}

func (ctr *RegisterController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", nil)
}

func (ctr *RegisterController) Register(c *gin.Context) {
	var request domain.RegisterUser

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	err = ctr.Validator.Validate(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	userData, err := request.ToUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	err = ctr.UserUsecase.Create(c, userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	c.JSON(http.StatusOK, domain.JsonResponse{
		Message: "Registration Successful",
		Success: true,
	})
}
