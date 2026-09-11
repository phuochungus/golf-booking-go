package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code: 0,
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, code int) {
	c.JSON(statusCode, ResponseData{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
