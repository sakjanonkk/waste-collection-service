package staff

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type staffService struct {
	repo domain.StaffReposity
}

func NewStaffService(repo domain.StaffReposity) domain.StaffService {
	return &staffService{repo: repo}
}

func (s *staffService) CreateStaff(staffInput models.Staff) (staff models.Staff, err error) {
	return s.repo.CreateStaff(staffInput)
}

func (s *staffService) GetStaffs(pagination models.Pagination) (staffs []models.Staff, paginated models.Pagination, err error) {
	return s.repo.GetStaffs(pagination)
}

func (s *staffService) GetStaffByID(staffInput models.Staff) (staff models.Staff, err error) {
	return s.repo.GetStaffByID(staffInput)
}

func (s *staffService) UpdateStaff(staffInput models.Staff) (staff models.Staff, err error) {
	return s.repo.UpdateStaff(staffInput)
}

func (s *staffService) DeleteStaff(staffInput models.Staff) error {
	return s.repo.DeleteStaff(staffInput)
}
