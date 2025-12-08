package vehicle

import (
	"errors"

	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) domain.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) CreateVehicle(vehicleInput models.Vehicle) (vehicle models.Vehicle, err error) {
	err = r.db.Create(&vehicleInput).Error
	if err != nil {
		return vehicle, err
	}
	return vehicleInput, nil
}

func (r *vehicleRepository) GetVehicles(pagination models.Pagination) (vehicles []models.Vehicle, paginated models.Pagination, err error) {
	var total int64
	err = r.db.Model(&models.Vehicle{}).Count(&total).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	paginated.Total = total

	err = r.db.Preload("Driver").
		Limit(pagination.PerPage).
		Offset((pagination.Page - 1) * pagination.PerPage).
		Find(&vehicles).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	return vehicles, paginated, nil
}

func (r *vehicleRepository) GetVehicleByID(vehicleInput models.Vehicle) (vehicle models.Vehicle, err error) {
	// ✅ เพิ่ม Preload เพื่อดึง driver info
	err = r.db.Preload("Driver").First(&vehicle, vehicleInput.ID).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (r *vehicleRepository) UpdateVehicle(vehicleInput models.Vehicle) (vehicle models.Vehicle, err error) {
	// ✅ Check ว่า record มีอยู่จริงหรือไม่
	if err := r.db.First(&vehicle, vehicleInput.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vehicle, errors.New("vehicle not found")
		}
		return vehicle, err
	}

	// ✅ Update with Select() เพื่อรองรับ NULL fields
	updates := map[string]interface{}{
		"registration_number":          vehicleInput.RegistrationNumber,
		"vehicle_type":                 vehicleInput.VehicleType,
		"vehicle_status":               vehicleInput.Status,
		"regular_waste_capacity_kg":    vehicleInput.RegularWasteCapacityKg,
		"recyclable_waste_capacity_kg": vehicleInput.RecyclableWasteCapacityKg,
		"fuel_type":                    vehicleInput.FuelType,
		"depreciation_value_per_year":  vehicleInput.DepreciationValuePerYear,
	}

	// ✅ จัดการ current_driver_id (รองรับการลบ driver)
	if vehicleInput.CurrentDriverID == nil {
		updates["current_driver_id"] = nil // set เป็น NULL
	} else if *vehicleInput.CurrentDriverID == 0 {
		updates["current_driver_id"] = nil // set เป็น NULL
	} else {
		updates["current_driver_id"] = *vehicleInput.CurrentDriverID
	}

	// ✅ จัดการ last_reported_problem
	if vehicleInput.LastReportedProblem == nil {
		updates["last_reported_problem"] = nil
	} else if *vehicleInput.LastReportedProblem == "" {
		updates["last_reported_problem"] = nil
	} else {
		updates["last_reported_problem"] = *vehicleInput.LastReportedProblem
	}

	err = r.db.Model(&vehicle).Updates(updates).Error
	if err != nil {
		return vehicle, err
	}

	// อ่านกลับมาอีกรอบพร้อม driver info
	err = r.db.Preload("Driver").First(&vehicle, vehicleInput.ID).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (r *vehicleRepository) DeleteVehicle(vehicleInput models.Vehicle) error {
	// ✅ Check ว่ามี record หรือไม่
	result := r.db.Delete(&models.Vehicle{}, vehicleInput.ID)
	if result.Error != nil {
		return result.Error
	}

	// ✅ ถ้าไม่มี record ที่ถูกลบ = ID ไม่มีอยู่
	if result.RowsAffected == 0 {
		return errors.New("vehicle not found")
	}

	return nil
}
