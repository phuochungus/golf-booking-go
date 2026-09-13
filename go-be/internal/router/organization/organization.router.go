package organization

import (
	"golf-booking-go/global"
	"golf-booking-go/internal/middlewares"
	"golf-booking-go/internal/modules/organization"

	"github.com/gin-gonic/gin"
)

func InitOrganizationRouter(apiGroup *gin.RouterGroup) {
	organizationService := organization.NewOrganizationService(global.DB)
	organizationController := organization.NewOrganizationController(organizationService)

	api := apiGroup.Group("/organization")
	api.Use(middlewares.AuthenticationMiddleware)
	{
		api.POST("/", organizationController.CreateOrganization)
		api.GET("/:id", organizationController.GetOrganizationByID)
	}
}
