package staff

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type staffReposity struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) domain.StaffReposity {
	return &staffReposity{db: db}
}

func (r *staffReposity) CreateStaff(staffInput models.Staff) (staff models.Staff, err error) {
	err = r.db.Create(&staffInput).Error
	if err != nil {
		return staff, err
	}
	return staffInput, nil
}

func (r *staffReposity) GetStaffs(pagination models.Pagination) (staffs []models.Staff, paginated models.Pagination, err error) {
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
