package response

import (
	"golf-booking-go/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details interface{} `json:"details,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code: 0,
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, code int, details interface{}) {
	if global.Config.Server.Mode == "dev" {
		c.JSON(statusCode, ResponseData{
			Code:    code,
			Message: message,
			Details: details,
		})
		return
	}
	c.JSON(statusCode, ResponseData{
		Code:    code,
		Message: message,
	})
}
