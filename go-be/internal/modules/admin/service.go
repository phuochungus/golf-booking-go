package admin

import (
	"errors"
	"fmt"
	"golf-booking-go/internal/entities"
	"golf-booking-go/internal/modules/admin/dto"
	"golf-booking-go/internal/utils"

	"gorm.io/gorm"
)

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

func (s *AdminService) RegisterAdmin(dto *dto.RegisterAdminDTO) (*int32, error) {
	var (
		admin    entities.Admin
		adminErr error
	)

	adminErr = s.db.
		Select("id").
		Where("email = ?", dto.Email).
		First(&admin).
		Error

	// Admin query
	switch {
	case adminErr == nil:
		return nil, errors.New("admin already exists")

	case !errors.Is(adminErr, gorm.ErrRecordNotFound):
		return nil, fmt.Errorf("check admin existence: %w", adminErr)
	}

	var returnedID *int32

	err := s.db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := utils.Hash(dto.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		admin := entities.Admin{
			Name:  dto.Email,
			Email: dto.Email,
		}

		if err := tx.Create(&admin).Error; err != nil {
			return fmt.Errorf("create admin: %w", err)
		}

		adminCred := entities.AdminCredential{
			AdminID:        admin.ID,
			HashedPassword: hashedPassword,
		}

		if err := tx.Create(&adminCred).Error; err != nil {
			return fmt.Errorf("create admin credential: %w", err)
		}

		returnedID = &admin.ID

		return nil
	})

	if err != nil {
		return nil, err
	}

	return returnedID, nil
}

func (s *AdminService) LoginAdmin(dto *dto.LoginAdminDTO) (string, string, error) {
	var (
		admin    entities.Admin
		adminErr error
	)

	adminErr = s.db.Joins("AdminCredential").
		Where("email = ?", dto.Email).
		First(&admin).
		Error

	switch {
	case errors.Is(adminErr, gorm.ErrRecordNotFound):
		return "", "", errors.New("admin not found")

	case adminErr != nil:
		return "", "", fmt.Errorf("check admin existence: %w", adminErr)
	}

	if err := utils.Compare(admin.AdminCredential.HashedPassword, dto.Password); err != nil {
		return "", "", errors.New("invalid password")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(admin.ID)
	if err != nil {
		return "", "", fmt.Errorf("generate JWT: %w", err)
	}

	return accessToken, refreshToken, nil
}
