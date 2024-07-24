package middleware

import "github.com/gin-gonic/gin"

func SetAuthCookies(accessToken string, refreshToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("accessToken", accessToken, 3600, "/", "localhost", true, true)
		c.SetCookie("refreshToken", refreshToken, 604800, "/", "localhost", true, true)
		c.Next()
	}
}
