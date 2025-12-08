package staff

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) domain.StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) CreateStaff(staffInput models.Staff) (staff models.Staff, err error) {
	err = r.db.Create(&staffInput).Error
	if err != nil {
		return staff, err
	}
	return staffInput, nil
}

func (r *staffRepository) GetStaffs(pagination models.Pagination) (staffs []models.Staff, paginated models.Pagination, err error) {
	var total int64
	err = r.db.Model(&models.Staff{}).Count(&total).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	paginated.Total = total
	err = r.db.Limit(pagination.PerPage).Offset((pagination.Page - 1) * pagination.PerPage).Find(&staffs).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	return staffs, paginated, nil
}

func (r *staffRepository) GetStaffByID(staffInput models.Staff) (staff models.Staff, err error) {
	err = r.db.First(&staff, staffInput.ID).Error
	if err != nil {
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) UpdateStaff(staffInput models.Staff) (staff models.Staff, err error) {
	err = r.db.Model(&models.Staff{}).Where("id = ?", staffInput.ID).Updates(staffInput).Error
	if err != nil {
		return staff, err
	}
	// อ่านกลับมาอีกรอบเพื่อคืนค่าเต็ม
	err = r.db.First(&staff, staffInput.ID).Error
	if err != nil {
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) DeleteStaff(staffInput models.Staff) error {
	return r.db.Delete(&models.Staff{}, staffInput.ID).Error
}

func (r *staffRepository) GetStaffByEmail(email string) (staff models.Staff, err error) {
	err = r.db.Where("email = ?", email).First(&staff).Error
	if err != nil {
		return staff, err
	}
	return staff, nil
}
