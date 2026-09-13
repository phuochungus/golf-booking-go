package organization

import (
	"context"
	"errors"
	"golf-booking-go/internal/entities"

	"gorm.io/gorm"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrAdminNotFound        = errors.New("admin not found")
	ErrAdminAlreadyInOrg    = errors.New("admin already belongs to an organization")
)

type OranizationService struct {
	db *gorm.DB
}

func NewOrganizationService(db *gorm.DB) *OranizationService {
	return &OranizationService{
		db: db,
	}
}

func (s *OranizationService) GetOrganizationByID(ctx context.Context, id int32) (*entities.Organization, error) {
	var organization entities.Organization
	err := s.db.WithContext(ctx).First(&organization, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrganizationNotFound
		}
		return nil, err
	}
	return &organization, nil
}

func (s *OranizationService) CreateOrganization(ctx context.Context, adminId, name string) (*int32, error) {
	var (
		admin        *entities.Admin
		organization *entities.Organization
	)
	s.db.WithContext(ctx).Model(&entities.Admin{}).Where("id = ?", adminId).First(&admin)

	switch {
	case admin == nil:
		return nil, ErrAdminNotFound
	case admin.OrganizationID != nil:
		return nil, ErrAdminAlreadyInOrg
	}

	organization = &entities.Organization{
		Name:        name,
		RootAdminID: &admin.ID,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(organization).Error; err != nil {
			return err
		}
		admin.OrganizationID = &organization.ID
		if err := tx.Save(admin).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &organization.ID, nil
}
