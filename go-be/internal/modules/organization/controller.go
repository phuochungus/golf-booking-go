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

func (oc *OrganizationController) CreateOrganization(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body", response.ErrorCodeBadRequest, err)
		return
	}

	adminId, exists := c.Get("adminId")

	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", response.ErrorCodeUnauthorized, nil)
		return
	}

	if adminId == nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", response.ErrorCodeUnauthorized, nil)
		return
	}

	adminIdInt := adminId.(int32)
	if adminIdInt == 0 {
		response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized", response.ErrorCodeUnauthorized, nil)
		return
	}

	organizationID, err := oc.organizationService.CreateOrganization(c.Request.Context(), adminIdInt, body.Name)
	switch err {
	case nil:
		response.SuccessResponse(c, gin.H{"organization_id": organizationID})
	case ErrAdminNotFound:
		response.ErrorResponse(c, http.StatusBadRequest, err.Error(), response.ErrorCodeAdminNotFound, nil)
	case ErrAdminAlreadyInOrg:
		response.ErrorResponse(c, http.StatusConflict, err.Error(), response.ErrorCodeAdminAlreadyInOrg, nil)
	case ErrOrganizationSameNameAlreadyExists:
		response.ErrorResponse(c, http.StatusConflict, err.Error(), response.ErrorCodeOrganizationSameNameAlreadyExists, nil)
	default:
		fmt.Println("Error creating organization:", err)
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error", response.ErrorCodeInternalServer, err)
	}
}
