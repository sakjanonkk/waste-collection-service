package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type StaffService interface {
	CreateStaff(models.Staff) (staff models.Staff, err error)
	GetStaffs(models.Pagination) (staffs []models.Staff, pagination models.Pagination, err error)
	GetStaffByID(models.Staff) (staff models.Staff, err error)
	UpdateStaff(models.Staff) (staff models.Staff, err error)
	DeleteStaff(models.Staff) error
	LoginStaff(email, password string) (staff models.Staff, err error)
}

type StaffRepository interface {
	CreateStaff(models.Staff) (staff models.Staff, err error)
	GetStaffs(models.Pagination) (staffs []models.Staff, pagination models.Pagination, err error)
	GetStaffByID(models.Staff) (staff models.Staff, err error)
	UpdateStaff(models.Staff) (staff models.Staff, err error)
	DeleteStaff(models.Staff) error
	GetStaffByEmail(email string) (staff models.Staff, err error)
}
