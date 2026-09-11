package admin

import (
	"errors"
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
	id, err := a.s.RegisterAdmin(c.Request.Context(), dto)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error(), response.ErrorCodeBadRequest)
		return
	}
	response.SuccessResponse(c, gin.H{"adminId": id})
}

func (a *AdminController) LoginAdmin(c *gin.Context) {
	var dto dto.LoginAdminDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, refreshToken, err := a.s.LoginAdmin(c.Request.Context(), &dto)
	if err != nil {
		status, message := http.StatusInternalServerError, "authentication unavailable"
		if errors.Is(err, ErrUnauthorized) {
			status, message = http.StatusUnauthorized, ErrUnauthorized.Error()
		} else {
			_ = c.Error(err)
		}
		response.ErrorResponse(c, status, message, status)
		return
	}
	c.Header("Cache-Control", "no-store")
	response.SuccessResponse(c, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

func (a *AdminController) RefreshAdmin(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required,max=4096"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "refreshToken is required and must be at most 4096 characters", http.StatusBadRequest)
		return
	}
	access, refresh, err := a.s.RefreshAdmin(c.Request.Context(), body.RefreshToken)
	if err != nil {
		status, message := http.StatusInternalServerError, "authentication unavailable"
		if errors.Is(err, ErrUnauthorized) {
			status, message = http.StatusUnauthorized, ErrUnauthorized.Error()
		} else {
			_ = c.Error(err)
		}
		response.ErrorResponse(c, status, message, status)
		return
	}
	c.Header("Cache-Control", "no-store")
	response.SuccessResponse(c, gin.H{"accessToken": access, "refreshToken": refresh})
}
