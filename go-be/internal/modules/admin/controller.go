package admin

import (
	"golf-booking-go/internal/modules/admin/dto"
	"golf-booking-go/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	s *AdminService
}

func NewAdminController(s *AdminService) *AdminController {
	return &AdminController{s: s}
}

func (a *AdminController) RegisterAdmin(c *gin.Context) {
	var dto *dto.RegisterAdminDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := a.s.RegisterAdmin(dto)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error(), response.ErrorCodeBadRequest)
		return
	}
	response.SuccessResponse(c, gin.H{"adminId": id})
}
