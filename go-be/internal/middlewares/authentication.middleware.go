package middlewares

import (
	"golf-booking-go/internal/utils"
	"golf-booking-go/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AuthenticationMiddleware gin.HandlerFunc
)

func InitAuthenticationMiddleware() {
	if AuthenticationMiddleware != nil {
		return
	}

	AuthenticationMiddleware = func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "missing access token", response.ErrorCodeUnauthorized, nil)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.ErrorResponse(c, http.StatusUnauthorized, "missing access token", response.ErrorCodeUnauthorized, nil)
			c.Abort()
			return
		}

		token := strings.Split(authHeader, " ")[1]

		claims, err := utils.ParseToken(token, "access")
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid access token", response.ErrorCodeUnauthorized, nil)
			c.Abort()
			return
		}
		c.Set("adminId", claims.AdminID)
		c.Next()
	}
}
