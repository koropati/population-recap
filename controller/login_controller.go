package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/bootstrap"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/internal/cryptos"
	"github.com/koropati/population-recap/internal/tokenutil"
	"github.com/koropati/population-recap/internal/validator"
	"github.com/koropati/population-recap/middleware"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	UserUsecase         domain.UserUsecase
	AccessTokenUsecase  domain.AccessTokenUsecase
	RefreshTokenUsecase domain.RefreshTokenUsecase
	Config              *bootstrap.Config
	Cryptos             cryptos.Cryptos
	Validator           *validator.Validator
}

func (ctr *LoginController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func (ctr *LoginController) Login(c *gin.Context) {
	var request domain.LoginUser

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

	user, err := ctr.UserUsecase.GetByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: "Wrong email or password", Success: false})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: "Wrong email or password", Success: false})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: "User is not active", Success: false})
		return
	}

	accessToken, err := tokenutil.CreateAccessToken(&user, ctr.Config.AccessTokenSecret, ctr.Config.AccessTokenExpiryHour, ctr.AccessTokenUsecase)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	refreshToken, err := tokenutil.CreateRefreshToken(&user, ctr.Config.RefreshTokenSecret, ctr.Config.RefreshTokenExpiryHour, accessToken, ctr.RefreshTokenUsecase)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	err = middleware.SetAuthContext(c, ctr.Cryptos, accessToken, refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	c.JSON(http.StatusOK, domain.JsonResponse{
		Message: "Login Successful",
		Success: true,
		Data: domain.UserTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
