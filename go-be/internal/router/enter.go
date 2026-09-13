package router

import (
	"golf-booking-go/internal/router/admin"
	"golf-booking-go/internal/router/organization"

	"github.com/gin-gonic/gin"
)

func InitRouter(apiGroup *gin.RouterGroup) {
	admin.InitAdminRouter(apiGroup)
	organization.InitOrganizationRouter(apiGroup)
}
