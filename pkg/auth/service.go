package auth

import (
	"errors"

	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	staffRepo    domain.StaffRepository
	jwtResources *models.JwtResources
}

func NewAuthService(staffRepo domain.StaffRepository, jwtResources *models.JwtResources) domain.AuthService {
	return &authService{
		staffRepo:    staffRepo,
		jwtResources: jwtResources,
	}
}

func (s *authService) Login(email, password string) (token string, staff models.Staff, err error) {
	// Get staff by email
	staff, err = s.staffRepo.GetStaffByEmail(email)
	if err != nil {
		return "", staff, errors.New("invalid email or password")
	}

	// Check if staff is active
	if staff.Status != models.StatusStaffActive {
		return "", staff, errors.New("account is not active")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(password))
	if err != nil {
		return "", staff, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err = GenerateToken(s.jwtResources, staff)
	if err != nil {
		return "", staff, errors.New("failed to generate token")
	}

	return token, staff, nil
}

func (s *authService) GetCurrentUser(staffID uint) (staff models.Staff, err error) {
	staff, err = s.staffRepo.GetStaffByID(models.Staff{ID: staffID})
	if err != nil {
		return staff, errors.New("user not found")
	}
	return staff, nil
}

func (s *authService) ChangePassword(staffID uint, oldPassword, newPassword string) error {
	// Get current user
	staff, err := s.staffRepo.GetStaffByID(models.Staff{ID: staffID})
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Update password
	staff.Password = string(hashedPassword)
	_, err = s.staffRepo.UpdateStaff(staff)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
