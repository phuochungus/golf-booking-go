package admin

import "golf-booking-go/internal/entities"

type AdminService struct {
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) RegisterAdmin(admin *entities.Admin) error {

}
