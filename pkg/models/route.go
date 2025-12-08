package models

import (
	"time"

	"gorm.io/gorm"
)

type RouteStatus string
type RoutePointStatus string

const (
	RouteStatusPlanned    RouteStatus = "planned"
	RouteStatusActive     RouteStatus = "active"
	RouteStatusComputeted RouteStatus = "completed"
	RouteStatusCancelled  RouteStatus = "cancelled"

	RoutePointStatusPending   RoutePointStatus = "pending"
	RoutePointStatusCollected RoutePointStatus = "collected"
	RoutePointStatusSkipped   RoutePointStatus = "skipped"
)

type Route struct {
	ID                   uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Date                 time.Time      `json:"date" gorm:"column:date;not null"`
	DriverID             uint           `json:"driver_id" gorm:"column:driver_id;not null"`
	Driver               *Staff         `json:"driver,omitempty" gorm:"foreignKey:DriverID"`
	VehicleID            uint           `json:"vehicle_id" gorm:"column:vehicle_id;not null"`
	Vehicle              *Vehicle       `json:"vehicle,omitempty" gorm:"foreignKey:VehicleID"`
	Status               RouteStatus    `json:"status" gorm:"column:status;not null;default:'planned'"`
	EstimatedDistance    float64        `json:"estimated_distance" gorm:"column:estimated_distance"`
	EstimatedTime        float64        `json:"estimated_time" gorm:"column:estimated_time"`
	FuelCostEstimate     float64        `json:"fuel_cost_estimate" gorm:"column:fuel_cost_estimate"`
	DepreciationEstimate float64        `json:"depreciation_estimate" gorm:"column:depreciation_estimate"`
	Notes                string         `json:"notes" gorm:"column:notes"`
	RoutePoints          []RoutePoint   `json:"route_points,omitempty" gorm:"foreignKey:RouteID"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

type RoutePoint struct {
	ID                 uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	RouteID            uint             `json:"route_id" gorm:"column:route_id;not null"`
	PointID            uint             `json:"point_id" gorm:"column:point_id;not null"`
	Point              *CollectionPoint `json:"point,omitempty" gorm:"foreignKey:PointID"`
	SequenceOrder      int              `json:"sequence_order" gorm:"column:sequence_order;not null"`
	CollectedAt        *time.Time       `json:"collected_at" gorm:"column:collected_at"`
	CollectedByID      *uint            `json:"collected_by_id" gorm:"column:collected_by_id"`
	CollectedBy        *Staff           `json:"collected_by,omitempty" gorm:"foreignKey:CollectedByID"`
	RegularWasteAmount float64          `json:"regular_waste_amount" gorm:"column:regular_waste_amount"`
	RecycleWasteAmount float64          `json:"recycle_waste_amount" gorm:"column:recycle_waste_amount"`
	WasteUnit          string           `json:"waste_unit" gorm:"column:waste_unit"`
	Status             RoutePointStatus `json:"status" gorm:"column:status;not null;default:'pending'"`
	Notes              string           `json:"notes" gorm:"column:notes"`
	CreatedAt          time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt   `json:"-" gorm:"index"`
}

type RouteInput struct {
	Date      string            `json:"date"` // YYYY-MM-DD
	DriverID  uint              `json:"driver_id"`
	VehicleID uint              `json:"vehicle_id"`
	Status    RouteStatus       `json:"status"`
	Points    []RoutePointInput `json:"points"`
	Notes     string            `json:"notes"`
}

type RoutePointInput struct {
	PointID       uint `json:"point_id"`
	SequenceOrder int  `json:"sequence_order"`
}

func (input *RouteInput) ToRoute() (Route, error) {
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return Route{}, err
	}
	return Route{
		Date:      parsedDate,
		DriverID:  input.DriverID,
		VehicleID: input.VehicleID,
		Status:    input.Status,
		Notes:     input.Notes,
	}, nil
}
