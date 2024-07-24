package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/bootstrap"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/internal/cryptos"
	"github.com/koropati/population-recap/internal/validator"
	"github.com/koropati/population-recap/middleware"
)

type LogoutController struct {
	UserUsecase         domain.UserUsecase
	AccessTokenUsecase  domain.AccessTokenUsecase
	RefreshTokenUsecase domain.RefreshTokenUsecase
	Config              *bootstrap.Config
	Cryptos             cryptos.Cryptos
	Validator           *validator.Validator
}

func (ctr *LogoutController) Logout(c *gin.Context) {

	refreshToken, errGetRefresh := middleware.GetAuthContext(c, ctr.Cryptos, "refresh")
	if errGetRefresh != nil {
		log.Printf("Error Get Refresh Token : %v\n", errGetRefresh)
	}
	accessToken, errGetAccess := middleware.GetAuthContext(c, ctr.Cryptos, "access")
	if errGetAccess != nil {
		log.Printf("Error Get Access Token : %v\n", errGetRefresh)
	}

	ctr.AccessTokenUsecase.Delete(c, accessToken)
	ctr.RefreshTokenUsecase.Delete(c, refreshToken)

	session := sessions.Default(c) // Get the current session

	// Clear the session data
	session.Clear()

	// Save the changes to the session store
	err := session.Save()
	if err != nil {
		// Handle error saving the session
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Redirect to the login page or other appropriate location
	c.Redirect(http.StatusFound, "/login")
}
