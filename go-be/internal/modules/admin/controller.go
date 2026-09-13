package admin

import (
	"errors"
	"golf-booking-go/internal/modules/admin/dto"
	"golf-booking-go/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService *AdminService
}

func NewAdminController(s *AdminService) *AdminController {
	return &AdminController{adminService: s}
}

func (a *AdminController) RegisterAdmin(c *gin.Context) {
	var dto *dto.RegisterAdminDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error(), http.StatusBadRequest, err)
		return
	}
	id, err := a.adminService.RegisterAdmin(c.Request.Context(), dto)
	switch err {
	case nil:
		response.SuccessResponse(c, gin.H{"adminId": id})
	case ErrAdminExists:
		response.ErrorResponse(c, http.StatusBadRequest, err.Error(), response.ErrorCodeAdminExisted, nil)
	default:
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error", response.ErrorCodeInternalServer, err)
	}
}

func (a *AdminController) LoginAdmin(c *gin.Context) {
	var dto dto.LoginAdminDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error(), http.StatusBadRequest, err)
		return
	}
	accessToken, refreshToken, err := a.adminService.LoginAdmin(c.Request.Context(), &dto)
	switch err {
	case nil:
		c.Header("Cache-Control", "no-store")
		response.SuccessResponse(c, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
	case ErrUnauthorized:
		response.ErrorResponse(c, http.StatusUnauthorized, err.Error(), response.ErrorCodeUnauthorized, nil)
	default:
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error", response.ErrorCodeInternalServer, err)
	}
}

func (a *AdminController) RefreshAdmin(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required,max=4096"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error(), http.StatusBadRequest, err)
		return
	}
	access, refresh, err := a.adminService.RefreshToken(c.Request.Context(), body.RefreshToken)
	if err != nil {
		status, message := http.StatusInternalServerError, "authentication unavailable"
		if errors.Is(err, ErrUnauthorized) {
			status, message = http.StatusUnauthorized, ErrUnauthorized.Error()
		} else {
			_ = c.Error(err)
		}
		response.ErrorResponse(c, status, message, status, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	response.SuccessResponse(c, gin.H{"accessToken": access, "refreshToken": refresh})
}
