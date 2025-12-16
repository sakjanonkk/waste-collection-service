package models

import "time"

type VehicleStatus string

const (
	StatusActive         VehicleStatus = "active"
	StatusInMaintenance  VehicleStatus = "in_maintenance"
	StatusDecommissioned VehicleStatus = "decommissioned"
)

type Vehicle struct {
	ID                        uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	RegistrationNumber        string        `json:"registration_number" gorm:"column:registration_number;not null;unique"`
	VehicleType               string        `json:"vehicle_type" gorm:"column:vehicle_type;not null"`
	Status                    VehicleStatus `json:"status" gorm:"column:vehicle_status;not null;default:'active'"`
	RegularWasteCapacityKg    float64       `json:"regular_waste_capacity_kg" gorm:"column:regular_waste_capacity_kg"`
	RecyclableWasteCapacityKg float64       `json:"recyclable_waste_capacity_kg" gorm:"column:recyclable_waste_capacity_kg"`
	CurrentDriverID           *uint         `json:"current_driver_id,omitempty" gorm:"column:current_driver_id"`
	Driver                    *Staff        `json:"driver,omitempty" gorm:"foreignKey:CurrentDriverID"`
	FuelType                  string        `json:"fuel_type" gorm:"column:fuel_type"`
	LastReportedProblem       *string       `json:"last_reported_problem,omitempty" gorm:"column:last_reported_problem"`
	DepreciationValuePerYear  float64       `json:"depreciation_value_per_year" gorm:"column:depreciation_value_per_year"`
	Image                     string        `json:"image" gorm:"column:image"`
	CreatedAt                 time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                 time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type VehicleInput struct {
	RegistrationNumber        string  `json:"registration_number" form:"registration_number"`
	VehicleType               string  `json:"vehicle_type" form:"vehicle_type"`
	Status                    string  `json:"status" form:"status"`
	RegularWasteCapacityKg    float64 `json:"regular_waste_capacity_kg" form:"regular_waste_capacity_kg"`
	RecyclableWasteCapacityKg float64 `json:"recyclable_waste_capacity_kg" form:"recyclable_waste_capacity_kg"`
	CurrentDriverID           *uint   `json:"current_driver_id,omitempty" form:"current_driver_id"`
	FuelType                  string  `json:"fuel_type" form:"fuel_type"`
	LastReportedProblem       *string `json:"last_reported_problem,omitempty" form:"last_reported_problem"`
	DepreciationValuePerYear  float64 `json:"depreciation_value_per_year" form:"depreciation_value_per_year"`
	Image                     string  `json:"image" form:"image"`
}

func (input *VehicleInput) ToVehicle() Vehicle {
	return Vehicle{
		RegistrationNumber:        input.RegistrationNumber,
		VehicleType:               input.VehicleType,
		Status:                    VehicleStatus(input.Status),
		RegularWasteCapacityKg:    input.RegularWasteCapacityKg,
		RecyclableWasteCapacityKg: input.RecyclableWasteCapacityKg,
		CurrentDriverID:           input.CurrentDriverID,
		FuelType:                  input.FuelType,
		LastReportedProblem:       input.LastReportedProblem,
		DepreciationValuePerYear:  input.DepreciationValuePerYear,
		Image:                     input.Image,
	}
}
