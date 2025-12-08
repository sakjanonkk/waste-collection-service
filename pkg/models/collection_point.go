package models

import (
	"time"

	"gorm.io/gorm"
)

type CollectionPointStatus string

const (
	StatusPointActive   CollectionPointStatus = "active"
	StatusPointInactive CollectionPointStatus = "inactive"
)

type CollectionPoint struct {
	ID              uint                  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string                `json:"name" gorm:"column:name;not null"`
	Latitude        float64               `json:"latitude" gorm:"column:latitude;not null"`
	Longitude       float64               `json:"longitude" gorm:"column:longitude;not null"`
	Address         string                `json:"address" gorm:"column:address"`
	Status          CollectionPointStatus `json:"status" gorm:"column:status;not null;default:'active'"`
	Image           string                `json:"image" gorm:"column:image"`
	ProblemReported string                `json:"problem_reported" gorm:"column:problem_reported"`
	RegularCapacity float64               `json:"regular_capacity" gorm:"column:regular_capacity"`
	RecycleCapacity float64               `json:"recycle_capacity" gorm:"column:recycle_capacity"`
	CreatedAt       time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt        `json:"-" gorm:"index"`
}

type CollectionPointInput struct {
	Name            string                `json:"name"`
	Latitude        float64               `json:"latitude"`
	Longitude       float64               `json:"longitude"`
	Address         string                `json:"address"`
	Status          CollectionPointStatus `json:"status"`
	Image           string                `json:"image"`
	ProblemReported string                `json:"problem_reported"`
	RegularCapacity float64               `json:"regular_capacity"`
	RecycleCapacity float64               `json:"recycle_capacity"`
}

func (input *CollectionPointInput) ToCollectionPoint() CollectionPoint {
	return CollectionPoint{
		Name:            input.Name,
		Latitude:        input.Latitude,
		Longitude:       input.Longitude,
		Address:         input.Address,
		Status:          input.Status,
		Image:           input.Image,
		ProblemReported: input.ProblemReported,
		RegularCapacity: input.RegularCapacity,
		RecycleCapacity: input.RecycleCapacity,
	}
}
