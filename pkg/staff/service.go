package staff

import (
	"errors"
	"strings"

	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type staffService struct {
	repo domain.StaffRepository
}

func NewStaffService(repo domain.StaffRepository) domain.StaffService {
	return &staffService{repo: repo}
}

func (s *staffService) CreateStaff(staffInput models.Staff) (staff models.Staff, err error) {
	if err := s.validateStaffInput(staffInput); err != nil {
		return staff, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(staffInput.Password),
		bcrypt.DefaultCost, // cost = 10
	)
	if err != nil {
		return staff, errors.New("failed to hash password")
	}
	staffInput.Password = string(hashedPassword)

	return s.repo.CreateStaff(staffInput)
}

func (s *staffService) validateStaffInput(staff models.Staff) error {
	if strings.TrimSpace(staff.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(staff.Password) == "" {
		return errors.New("password is required")
	}
	if len(staff.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if strings.TrimSpace(staff.FirstName) == "" {
		return errors.New("first name is required")
	}
	if strings.TrimSpace(staff.LastName) == "" {
		return errors.New("last name is required")
	}

	validRoles := map[models.StaffRole]bool{
		models.RoleAdmin:        true,
		models.RoleRouteManager: true,
		models.RoleDriver:       true,
		models.RoleCollector:    true,
		models.RoleCitizen:      true,
	}
	if !validRoles[staff.Role] {
		return errors.New("invalid role")
	}

	validStatuses := map[models.StaffStatus]bool{
		models.StatusStaffActive:   true,
		models.StatusStaffInactive: true,
		models.StatusStaffOnLeave:  true,
	}
	if !validStatuses[staff.Status] {
		return errors.New("invalid status")
	}

	return nil
}

func (s *staffService) GetStaffs(pagination models.Pagination, filter models.StaffFilter) (staffs []models.Staff, paginated models.Pagination, err error) {
	return s.repo.GetStaffs(pagination, filter)
}

func (s *staffService) GetStaffByID(staffInput models.Staff) (staff models.Staff, err error) {
	return s.repo.GetStaffByID(staffInput)
}

func (s *staffService) UpdateStaff(staffInput models.Staff) (staff models.Staff, err error) {
	if strings.TrimSpace(staffInput.Password) != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(staffInput.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return staff, errors.New("failed to hash password")
		}
		staffInput.Password = string(hashedPassword)
	}

	return s.repo.UpdateStaff(staffInput)
}

func (s *staffService) DeleteStaff(staffInput models.Staff) error {
	return s.repo.DeleteStaff(staffInput)
}

func (s *staffService) LoginStaff(email, password string) (staff models.Staff, err error) {
	staff, err = s.repo.GetStaffByEmail(email)
	if err != nil {
		return staff, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(staff.Password),
		[]byte(password),
	)
	if err != nil {
		return staff, errors.New("invalid email or password")
	}

	return staff, nil
}
