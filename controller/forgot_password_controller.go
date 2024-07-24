package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/bootstrap"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/internal/cryptos"
	"github.com/koropati/population-recap/internal/mailer"
	"github.com/koropati/population-recap/internal/tokenutil"
	"github.com/koropati/population-recap/internal/urlutil"
	"github.com/koropati/population-recap/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordController struct {
	Config                     *bootstrap.Config
	Cryptos                    cryptos.Cryptos
	Validator                  *validator.Validator
	UserUsecase                domain.UserUsecase
	ForgotPasswordTokenUsecase domain.ForgotPasswordTokenUsecase
	Mailer                     mailer.Mailer
}

const (
	ForgotPasswordPath       = "/forgot-password"
	VerifyForgotPasswordPath = "/forgot-password/verify"
	ResetPasswordPath        = "/reset-password"
)

func (ctr *ForgotPasswordController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "forgot_password.tmpl", nil)
}

func (ctr *ForgotPasswordController) ForgotPassword(c *gin.Context) {
	var request domain.ForgotPassword

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
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: "Email not registered", Success: false})
		return
	}

	forgotToken, err := tokenutil.CreateForgotToken(&user, ctr.Config.ForgotTokenExpiryHour, ctr.ForgotPasswordTokenUsecase)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	encForgotPasswordToken, err := ctr.Cryptos.Encrypt(forgotToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	go func() {
		errSendEmail := ctr.Mailer.SendForgotPassword(mailer.ForgotPasswordData{
			AppName:             ctr.Config.AppName,
			Name:                user.Name,
			Email:               user.Email,
			From:                ctr.Config.SmtpSenderMail,
			ForgotPasswordToken: forgotToken,
			UrlRedirect:         urlutil.CreateUrlForgotPassword(c.Request, encForgotPasswordToken),
			UrlVerification:     urlutil.CreateUrlForgotPassword(c.Request, encForgotPasswordToken),
		})
		if errSendEmail != nil {
			log.Printf("Error Send Email : %v\n", errSendEmail)
		}
	}()

	c.JSON(http.StatusOK, domain.JsonResponse{
		Message: "Forgpt Password Link Send Successful",
		Success: true,
	})
}

func (ctr *ForgotPasswordController) VerifyForgotPassword(c *gin.Context) {
	var isSuccess bool
	var msg string
	var redirectURL string

	params := c.Request.URL.Query()
	token := params.Get("token")

	if token == "" {
		isSuccess = false
		msg = "Error when try to verification forgot password link, token data is invalid"

		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, msg, token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)
		return
	}

	token, err := ctr.Cryptos.Decrypt(token)
	if err != nil {
		isSuccess = false
		msg = "Error when try to verification forgot password link, invalid token"

		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, msg, token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)
		return
	}

	isValid := ctr.ForgotPasswordTokenUsecase.IsValid(c, token)
	if !isValid {
		isSuccess = false
		msg = "Error when try to verification forgot password link, invalid token"

		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, msg, token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)

		return
	}

	userId, err := ctr.ForgotPasswordTokenUsecase.GetUserID(c, token)
	if err != nil {
		isSuccess = false
		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, err.Error(), token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}

	user, err := ctr.UserUsecase.GetById(c, userId)
	if err != nil {
		isSuccess = false
		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, err.Error(), token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}

	err = ctr.ForgotPasswordTokenUsecase.Revoke(c, token)
	if err != nil {
		isSuccess = false
		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, err.Error(), token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}

	forgotToken, err := tokenutil.CreateForgotToken(&user, ctr.Config.ForgotTokenExpiryHour, ctr.ForgotPasswordTokenUsecase)
	if err != nil {
		isSuccess = false
		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, err.Error(), token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}

	encForgotPasswordToken, err := ctr.Cryptos.Encrypt(forgotToken)
	if err != nil {
		isSuccess = false
		redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ForgotPasswordPath, isSuccess, err.Error(), token, "")
		c.Redirect(http.StatusMovedPermanently, redirectURL)
	}

	isSuccess = true
	msg = "Success Validate Forgot Password Token"

	redirectURL = urlutil.GenerateRedirectForgotPassword(ctr.Config.AppFeUrl, ResetPasswordPath, isSuccess, msg, encForgotPasswordToken, user.Name)
	c.Redirect(http.StatusMovedPermanently, redirectURL)
}

func (ctr *ForgotPasswordController) ResetPassword(c *gin.Context) {
	params := c.Request.URL.Query()
	token := params.Get("t")
	success := params.Get("success")
	name := params.Get("name")
	msg := params.Get("msg")

	data := map[string]interface{}{
		"success": success,
		"msg":     msg,
		"token":   token,
		"name":    name,
	}
	c.HTML(http.StatusOK, "reset_password.tmpl", data)
}

func (ctr *ForgotPasswordController) ResetPasswordConfirm(c *gin.Context) {
	var request domain.ResetPassword

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

	if request.Password != request.RePassword {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: "Password not match", Success: false})
		return
	}

	token, err := ctr.Cryptos.Decrypt(request.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	isValid := ctr.ForgotPasswordTokenUsecase.IsValid(c, token)
	if !isValid {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: "Token is Expired", Success: false})
		return
	}

	userId, err := ctr.ForgotPasswordTokenUsecase.GetUserID(c, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	user, err := ctr.UserUsecase.GetById(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	err = ctr.ForgotPasswordTokenUsecase.Revoke(c, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	err = ctr.UserUsecase.UpdatePassword(c, user.ID, string(newPasswordHash))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Error(), Success: false})
		return
	}

	go func() {
		errSendEmail := ctr.Mailer.SendNotification(mailer.Notification{
			AppName: ctr.Config.AppName,
			Name:    user.Name,
			Email:   user.Email,
			From:    ctr.Config.SmtpSenderMail,
			Title:   "Password has been reset",
			Subject: "Reset Password Success",
			Message: "Your password has been successfuly reset, please login with your new password, have a good day :).",
		})
		if errSendEmail != nil {
			log.Printf("Error Send Email : %v\n", errSendEmail)
		}
	}()

	c.JSON(http.StatusOK, domain.JsonResponse{
		Message: "Your password has been reset :)",
		Success: true,
	})
}
