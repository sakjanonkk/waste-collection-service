package staff

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/utils"
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
	var existingRole models.Role
	if err := r.db.Where("name = ?", utils.User).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Staff{}, err
		}
		return models.Staff{}, err
	}
	var existingUserRole models.UserRole
	if err := r.db.Where("user_id = ? AND role_id = ?", staffInput.ID, existingRole.ID).First(&existingUserRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userRole := models.UserRole{
				UserID: staffInput.ID,
				RoleID: existingRole.ID,
			}
			if err := r.db.Create(&userRole).Error; err != nil {
				return models.Staff{}, err
			}
		} else {
			return models.Staff{}, err
		}
	}
	return staffInput, nil
}

func (r *staffRepository) GetStaffs(pagination models.Pagination, filter models.StaffFilter) (staffs []models.Staff, paginated models.Pagination, err error) {
	query := r.db.Model(&models.Staff{})

	// Apply filters
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("firstname LIKE ? OR lastname LIKE ? OR email LIKE ?", searchTerm, searchTerm, searchTerm)
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	paginated.Total = total
	err = query.Limit(pagination.PerPage).Offset((pagination.Page - 1) * pagination.PerPage).Find(&staffs).Error
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
