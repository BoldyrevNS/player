package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shared/common/response"
	"shared/common/token"
)

func AdminPermissionMiddleware(c *gin.Context) {
	headerToken, err := token.GetTokenFromHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.MessageJSON{
			Message: err.Error(),
		})
		return
	}
	accessTokenClaims, err := token.ParseAccessToken(headerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.MessageJSON{
			Message: err.Error(),
		})
		return
	}
	if accessTokenClaims.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.MessageJSON{
			Message: "permission denied",
		})
		return
	}
}

func ProtectedMiddleware(c *gin.Context) {
	headerToken, err := token.GetTokenFromHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.MessageJSON{
			Message: err.Error(),
		})
		return
	}
	_, err = token.ParseAccessToken(headerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.MessageJSON{
			Message: err.Error(),
		})
		return
	}
}
