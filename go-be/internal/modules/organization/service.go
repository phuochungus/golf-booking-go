package organization

import (
	"context"
	"errors"
	"golf-booking-go/internal/entities"
	"sync"

	"gorm.io/gorm"
)

var (
	ErrOrganizationNotFound              = errors.New("organization not found")
	ErrAdminNotFound                     = errors.New("admin not found")
	ErrAdminAlreadyInOrg                 = errors.New("admin already belongs to an organization")
	ErrOrganizationSameNameAlreadyExists = errors.New("organization with the same name already exists")
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

func (s *OranizationService) CreateOrganization(
	ctx context.Context,
	adminID int32,
	name string,
) (*int32, error) {
	var (
		admin        entities.Admin
		organization entities.Organization
		adminErr     error
		orgErr       error
		wg           sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		adminErr = s.db.WithContext(ctx).First(&admin, adminID).Error
	}()

	go func() {
		defer wg.Done()
		orgErr = s.db.WithContext(ctx).
			Where("name = ?", name).
			First(&organization).Error
	}()

	wg.Wait()

	if adminErr != nil {
		if errors.Is(adminErr, gorm.ErrRecordNotFound) {
			return nil, ErrAdminNotFound
		}
		return nil, adminErr
	}

	if admin.OrganizationID != nil {
		return nil, ErrAdminAlreadyInOrg
	}

	if orgErr == nil {
		return nil, ErrOrganizationSameNameAlreadyExists
	}
	if !errors.Is(orgErr, gorm.ErrRecordNotFound) {
		return nil, orgErr
	}

	var organizationID int32

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Re-check admin inside transaction to avoid concurrent requests
		// assigning the same admin to multiple organizations.
		var lockedAdmin entities.Admin
		if err := tx.First(&lockedAdmin, adminID).Error; err != nil {
			return err
		}
		if lockedAdmin.OrganizationID != nil {
			return ErrAdminAlreadyInOrg
		}

		newOrganization := entities.Organization{
			Name:        name,
			RootAdminID: &lockedAdmin.ID,
		}
		if err := tx.Create(&newOrganization).Error; err != nil {
			return err
		}

		if err := tx.Model(&lockedAdmin).
			Update("organization_id", newOrganization.ID).Error; err != nil {
			return err
		}

		organizationID = newOrganization.ID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &organizationID, nil
}
