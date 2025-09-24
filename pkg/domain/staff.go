package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type StaffService interface {
	CreateStaff(models.Staff) (staff models.Staff, err error)
	GetStaffs(models.Pagination) (staffs []models.Staff, pagination models.Pagination, err error)
	GetStaffByID(models.Staff) (staff models.Staff, err error)
	UpdateStaff(models.Staff) (staff models.Staff, err error) // ✅ update
	DeleteStaff(models.Staff) error                           // ✅ delete
}
type StaffReposity interface {
	CreateStaff(models.Staff) (staff models.Staff, err error)
	GetStaffs(models.Pagination) (staffs []models.Staff, pagination models.Pagination, err error)
	GetStaffByID(models.Staff) (staff models.Staff, err error)
	UpdateStaff(models.Staff) (staff models.Staff, err error) // ✅ update
	DeleteStaff(models.Staff) error                           // ✅ delete
}
