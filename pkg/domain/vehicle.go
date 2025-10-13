package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type VehicleService interface {
	CreateVehicle(models.Vehicle) (vehicle models.Vehicle, err error)
	GetVehicles(models.Pagination) (vehicles []models.Vehicle, pagination models.Pagination, err error)
	GetVehicleByID(models.Vehicle) (vehicle models.Vehicle, err error)
	UpdateVehicle(models.Vehicle) (vehicle models.Vehicle, err error)
	DeleteVehicle(models.Vehicle) error
}

type VehicleRepository interface {
	CreateVehicle(models.Vehicle) (vehicle models.Vehicle, err error)
	GetVehicles(models.Pagination) (vehicles []models.Vehicle, pagination models.Pagination, err error)
	GetVehicleByID(models.Vehicle) (vehicle models.Vehicle, err error)
	UpdateVehicle(models.Vehicle) (vehicle models.Vehicle, err error)
	DeleteVehicle(models.Vehicle) error
}
