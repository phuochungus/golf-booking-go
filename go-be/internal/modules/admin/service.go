package admin

import (
	"errors"
	"fmt"
	"golf-booking-go/internal/entities"
	"golf-booking-go/internal/modules/admin/dto"
	"golf-booking-go/internal/utils"
	"sync"

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
		org      entities.Organization
		adminErr error
		orgErr   error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		adminErr = s.db.
			Select("id").
			Where("email = ?", dto.Email).
			First(&admin).
			Error
	}()

	go func() {
		defer wg.Done()

		orgErr = s.db.
			Select("id").
			Where("name = ?", dto.OrganizationName).
			First(&org).
			Error
	}()

	wg.Wait()

	// Admin query
	switch {
	case adminErr == nil:
		return nil, errors.New("admin already exists")

	case !errors.Is(adminErr, gorm.ErrRecordNotFound):
		return nil, fmt.Errorf("check admin existence: %w", adminErr)
	}

	// Organization query
	switch {
	case orgErr == nil:
		return nil, errors.New("organization already exists")

	case !errors.Is(orgErr, gorm.ErrRecordNotFound):
		return nil, fmt.Errorf("check organization existence: %w", orgErr)
	}

	var returnedID *int32

	err := s.db.Transaction(func(tx *gorm.DB) error {
		org := entities.Organization{
			Name: dto.OrganizationName,
		}

		if err := tx.Create(&org).Error; err != nil {
			return fmt.Errorf("create organization: %w", err)
		}

		hashedPassword, err := utils.Hash(dto.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		admin := entities.Admin{
			Name:           dto.Email,
			Email:          dto.Email,
			OrganizationID: &org.ID,
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
