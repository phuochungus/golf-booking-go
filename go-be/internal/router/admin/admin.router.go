package admin

import (
	"golf-booking-go/global"
	"golf-booking-go/internal/modules/admin"

	"github.com/gin-gonic/gin"
)

func InitAdminRouter(apiGroup *gin.RouterGroup) {
	service := admin.NewAdminService(global.DB)
	controller := admin.NewAdminController(service)

	adminGroup := apiGroup.Group("/admin")

	adminGroup.POST("/register", controller.RegisterAdmin)
	adminGroup.POST("/login", controller.LoginAdmin)
	adminGroup.POST("/refresh", controller.RefreshAdmin)
}
