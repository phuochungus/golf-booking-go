package admin

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"golf-booking-go/global"
	"golf-booking-go/internal/entities"
	"golf-booking-go/internal/modules/admin/dto"
	"golf-booking-go/internal/utils"

	"gorm.io/gorm"
)

var ErrUnauthorized = errors.New("invalid credentials or refresh token")

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

func (s *AdminService) RegisterAdmin(ctx context.Context, dto *dto.RegisterAdminDTO) (*int32, error) {
	var (
		admin    entities.Admin
		adminErr error
	)

	adminErr = s.db.WithContext(ctx).
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

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (s *AdminService) LoginAdmin(ctx context.Context, dto *dto.LoginAdminDTO) (string, string, error) {
	var admin entities.Admin

	adminErr := s.db.WithContext(ctx).Joins("AdminCredential").
		Where("email = ?", dto.Email).
		First(&admin).
		Error

	switch {
	case errors.Is(adminErr, gorm.ErrRecordNotFound):
		return "", "", ErrUnauthorized

	case adminErr != nil:
		return "", "", fmt.Errorf("check admin existence: %w", adminErr)
	}

	if admin.AdminCredential == nil {
		return "", "", ErrUnauthorized
	}
	if err := utils.Compare(admin.AdminCredential.HashedPassword, dto.Password); err != nil {
		return "", "", ErrUnauthorized
	}

	return s.issueTokens(ctx, admin.ID, "")
}

func (s *AdminService) RefreshAdmin(ctx context.Context, raw string) (string, string, error) {
	claims, err := utils.ParseToken(raw, "refresh")
	if errors.Is(err, utils.ErrInvalidToken) {
		return "", "", ErrUnauthorized
	}
	if err != nil {
		return "", "", err
	}
	var admin entities.Admin
	err = s.db.WithContext(ctx).Select("id").First(&admin, claims.AdminID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", ErrUnauthorized
	}
	if err != nil {
		return "", "", fmt.Errorf("find admin: %w", err)
	}
	return s.issueTokens(ctx, admin.ID, raw)
}

func (s *AdminService) issueTokens(ctx context.Context, adminID int32, previous string) (string, string, error) {
	if global.RDB == nil {
		return "", "", errors.New("refresh token store unavailable")
	}
	access, refresh, err := utils.GenerateTokens(adminID)
	if err != nil {
		return "", "", fmt.Errorf("generate JWT: %w", err)
	}
	key := fmt.Sprintf("admin:refresh:%x", sha256.Sum256([]byte(refresh)))
	if previous == "" {
		err = global.RDB.Set(ctx, key, "1", utils.RefreshTokenTTL).Err()
	} else {
		var rotated int
		// Only one concurrent request can replace this refresh token.
		previousKey := fmt.Sprintf("admin:refresh:%x", sha256.Sum256([]byte(previous)))
		rotated, err = global.RDB.Eval(ctx, `
if redis.call('EXISTS', KEYS[1]) == 0 then return 0 end
redis.call('SET', KEYS[2], '1', 'EX', ARGV[1])
redis.call('DEL', KEYS[1])
return 1
`, []string{previousKey, key}, int64(utils.RefreshTokenTTL.Seconds())).Int()
		if err == nil && rotated != 1 {
			return "", "", ErrUnauthorized
		}
	}
	if err != nil {
		return "", "", fmt.Errorf("store refresh token: %w", err)
	}
	return access, refresh, nil
}
