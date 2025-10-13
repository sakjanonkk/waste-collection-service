package vehicle

import (
	"errors"
	"strings"

	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type vehicleService struct {
	repo domain.VehicleRepository
}

func NewVehicleService(repo domain.VehicleRepository) domain.VehicleService {
	return &vehicleService{repo: repo}
}

func (s *vehicleService) CreateVehicle(vehicleInput models.Vehicle) (vehicle models.Vehicle, err error) {
	// ✅ เพิ่ม Validation
	if err := s.validateVehicle(vehicleInput); err != nil {
		return vehicle, err
	}

	return s.repo.CreateVehicle(vehicleInput)
}

// ✅ Validation function
func (s *vehicleService) validateVehicle(vehicle models.Vehicle) error {
	// Required fields
	if strings.TrimSpace(vehicle.RegistrationNumber) == "" {
		return errors.New("registration number is required")
	}
	if strings.TrimSpace(vehicle.VehicleType) == "" {
		return errors.New("vehicle type is required")
	}
	if strings.TrimSpace(vehicle.FuelType) == "" {
		return errors.New("fuel type is required")
	}

	// Validate status
	validStatuses := map[models.VehicleStatus]bool{
		models.StatusActive:         true,
		models.StatusInMaintenance:  true,
		models.StatusDecommissioned: true,
	}
	if !validStatuses[vehicle.Status] {
		return errors.New("invalid status: must be 'active', 'in_maintenance', or 'decommissioned'")
	}

	// Validate capacity
	if vehicle.RegularWasteCapacityKg < 0 {
		return errors.New("regular waste capacity cannot be negative")
	}
	if vehicle.RecyclableWasteCapacityKg < 0 {
		return errors.New("recyclable waste capacity cannot be negative")
	}
	if vehicle.DepreciationValuePerYear < 0 {
		return errors.New("depreciation value cannot be negative")
	}

	return nil
}

func (s *vehicleService) GetVehicles(pagination models.Pagination) (vehicles []models.Vehicle, paginated models.Pagination, err error) {
	return s.repo.GetVehicles(pagination)
}

func (s *vehicleService) GetVehicleByID(vehicleInput models.Vehicle) (vehicle models.Vehicle, err error) {
	return s.repo.GetVehicleByID(vehicleInput)
}

func (s *vehicleService) UpdateVehicle(vehicleInput models.Vehicle) (vehicle models.Vehicle, err error) {
	if err := s.validateVehicle(vehicleInput); err != nil {
		return vehicle, err
	}

	return s.repo.UpdateVehicle(vehicleInput)
}

func (s *vehicleService) DeleteVehicle(vehicleInput models.Vehicle) error {
	return s.repo.DeleteVehicle(vehicleInput)
}
