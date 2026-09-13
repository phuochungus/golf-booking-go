package organization

import (
	"fmt"
	"golf-booking-go/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	organizationService *OranizationService
}

func NewOrganizationController(organizationService *OranizationService) *OrganizationController {
	return &OrganizationController{
		organizationService: organizationService,
	}
}

func (oc *OrganizationController) GetOrganizationByID(c *gin.Context) {
	idParam := c.Param("id")
	var id int32
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		response.ErrorResponse(c, response.ErrorCodeBadRequest, "invalid organization ID", http.StatusBadRequest, err)
		return
	}

	organization, err := oc.organizationService.GetOrganizationByID(c.Request.Context(), id)
	switch err {
	case nil:
		response.SuccessResponse(c, gin.H{"organization": organization})
	case ErrOrganizationNotFound:
		response.ErrorResponse(c, response.ErrorCodeBadRequest, "organization not found", http.StatusNotFound, nil)
	default:
		response.ErrorResponse(c, response.ErrorCodeInternalServer, "internal server error", http.StatusInternalServerError, err)
	}
}

func (oc *OrganizationController) CreateOrganization(c *gin.Context, adminId, name string) {
	organizationID, err := oc.organizationService.CreateOrganization(c.Request.Context(), adminId, name)
	switch err {
	case nil:
		response.SuccessResponse(c, gin.H{"organization_id": organizationID})
	case ErrAdminNotFound:
		response.ErrorResponse(c, response.ErrorCodeBadRequest, err.Error(), http.StatusNotFound, nil)
	case ErrAdminAlreadyInOrg:
		response.ErrorResponse(c, response.ErrorCodeBadRequest, err.Error(), http.StatusConflict, nil)
	default:
		response.ErrorResponse(c, response.ErrorCodeInternalServer, "internal server error", http.StatusInternalServerError, err)
	}
}
