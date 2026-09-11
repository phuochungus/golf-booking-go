package dto

type RegisterAdminDTO struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required"`
	OrganizationName string `json:"organization_name" binding:"required"`
}
