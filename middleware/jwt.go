package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koropati/population-recap/domain"
	"github.com/koropati/population-recap/internal/cryptos"
	randomstr "github.com/koropati/population-recap/internal/reandomstr"
	"github.com/koropati/population-recap/internal/tokenutil"
	"github.com/koropati/population-recap/internal/urlutil"
)

const (
	UserIDContext      = "x-user-id"
	UserRoleContext    = "x-user-role"
	AuthAccessContext  = "x-a-auth"
	AuthRefreshContext = "x-r-auth"
	RoleSuperAdmin     = "super_admin"
	RoleAdmin          = "admin"
	RoleStaff          = "staff"
	RoleAnonymous      = "anonymous"
	RefreshToken       = "refresh_token"
	AccessToken        = "access_token"
)

func JwtAuthMiddleware(secret string, casbinEnforcer *casbin.Enforcer, cryptos cryptos.Cryptos, accessTokenUsecase domain.AccessTokenUsecase, refreshTokenUsecase domain.RefreshTokenUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		authToken, err := parseAuthorizationHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.JsonResponse{Message: err.Error(), Success: false})
			c.Abort()
			return
		}

		if !accessTokenUsecase.IsValid(c, authToken) {
			c.JSON(http.StatusUnauthorized, domain.JsonResponse{Message: "Token tidak valid atau telah kadaluarsa", Success: false})
			c.Abort()
			return
		}

		userID, userRole, err := tokenutil.ExtractIDFromToken(authToken, secret, AccessToken, accessTokenUsecase, refreshTokenUsecase)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.JsonResponse{Message: err.Error(), Success: false})
			c.Abort()
			return
		}

		if userRole == "" {
			userRole = RoleAnonymous
		}

		SetUserContext(c, cryptos, userID, userRole)

		if err := enforceCasbinRules(c, casbinEnforcer, userRole); err != nil {
			c.JSON(http.StatusBadRequest, domain.JsonResponse{Message: err.Message, Success: false})
			c.Abort()
			return
		}

		c.Next()
	}
}

func parseAuthorizationHeader(authHeader string) (string, error) {
	t := strings.Split(authHeader, " ")
	if len(t) != 2 {
		return "", errors.New("not authorized")
	}
	return t[1], nil
}

func SetAuthContext(c *gin.Context, cryptos cryptos.Cryptos, accessToken string, refreshToken string) error {
	encryptedAccessToken, err := cryptos.Encrypt(accessToken)
	if err != nil {
		return err
	}
	encryptedRefreshToken, err := cryptos.Encrypt(refreshToken)
	if err != nil {
		return err
	}

	c.Set(AuthAccessContext, encryptedAccessToken)
	c.Set(AuthRefreshContext, encryptedRefreshToken)

	session := sessions.Default(c)
	session.Set(AuthAccessContext, encryptedAccessToken)
	session.Set(AuthRefreshContext, encryptedRefreshToken)
	session.Save()
	return nil
}

func SetUserContext(c *gin.Context, cryptos cryptos.Cryptos, userID, userRole string) {
	encryptedUserID, err := cryptos.Encrypt(userID)
	if err != nil {
		panic(err)
	}

	randomStr := randomstr.New(true, true, true, false)

	garbage := randomStr.GenerateRandomString(5)

	userRoleWithGarbage := userRole + "_" + garbage

	encryptedUserRole, err := cryptos.Encrypt(userRoleWithGarbage)
	if err != nil {
		panic(err)
	}
	c.Set(UserIDContext, encryptedUserID)
	c.Set(UserRoleContext, encryptedUserRole)
}

func GetUserContext(c *gin.Context, cryptos cryptos.Cryptos) (userID string, userRole string) {
	encryptedUserID := c.GetString(UserIDContext)
	encryptedUserRole := c.GetString(UserRoleContext)

	if encryptedUserID == "" || encryptedUserRole == "" {
		return "", ""
	}
	userID, err := cryptos.Decrypt(encryptedUserID)
	if err != nil {
		panic(err)
	}
	userRoleWithGarbage, err := cryptos.Decrypt(encryptedUserRole)
	if err != nil {
		panic(err)
	}
	userRole = strings.Split(userRoleWithGarbage, "_")[0]
	return userID, userRole
}

func GetAuthContext(c *gin.Context, cryptos cryptos.Cryptos, tokenType string) (token string, err error) {
	encryptedToken := ""
	var dataSession interface{}
	session := sessions.Default(c)
	if tokenType == "refresh" {
		dataSession = session.Get(AuthRefreshContext)
	} else {
		dataSession = session.Get(AuthAccessContext)
	}
	if dataSession == nil {
		return "", errors.New("invalid session")
	}

	encryptedToken = dataSession.(string)

	token, err = cryptos.Decrypt(encryptedToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

func enforceCasbinRules(c *gin.Context, casbinEnforcer *casbin.Enforcer, userRole string) *domain.JsonResponse {
	pathUrl := urlutil.RemoveAPIVersionMiddleware(c.Request.URL.Path)
	res, err := casbinEnforcer.EnforceSafe(userRole, pathUrl, c.Request.Method)
	if err != nil {
		return &domain.JsonResponse{Message: err.Error(), Success: false}
	}
	if !res {
		return &domain.JsonResponse{Message: "unauthorized", Success: false}
	}
	return nil
}
