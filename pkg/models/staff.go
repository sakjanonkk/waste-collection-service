package models

import (
	"time"

	"gorm.io/gorm"
)

type StaffRole string

const (
	RoleAdmin        StaffRole = "admin"
	RoleRouteManager StaffRole = "route_manager"
	RoleDriver       StaffRole = "driver"
	RoleCollector    StaffRole = "collector"
	RoleCitizen      StaffRole = "citizen"
)

type StaffStatus string

const (
	StatusStaffActive   StaffStatus = "active"
	StatusStaffInactive StaffStatus = "inactive"
	StatusStaffOnLeave  StaffStatus = "on_leave"
)

type Staff struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix      string         `json:"prefix" gorm:"column:prefix;not null"`
	FirstName   string         `json:"firstname" gorm:"column:firstname;not null"`
	LastName    string         `json:"lastname" gorm:"column:lastname;not null"`
	Email       string         `json:"email" gorm:"column:email;not null;unique"`
	Password    string         `json:"-" gorm:"column:password;not null"`
	Role        StaffRole      `json:"role" gorm:"column:role;not null"`
	Status      StaffStatus    `json:"status" gorm:"column:status;not null"`
	PhoneNumber string         `json:"phone_number" gorm:"column:phone;not null"`
	Picture     string         `json:"picture" gorm:"column:picture"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type StaffInput struct {
	Prefix      string      `json:"prefix"`
	FirstName   string      `json:"firstname"`
	LastName    string      `json:"lastname"`
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	Role        StaffRole   `json:"role"`
	Status      StaffStatus `json:"status"`
	PhoneNumber string      `json:"phone_number"`
	Picture     string      `json:"picture"`
}

func (input *StaffInput) ToStaff() Staff {
	return Staff{
		Prefix:      input.Prefix,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Password:    input.Password,
		Role:        input.Role,
		Status:      input.Status,
		PhoneNumber: input.PhoneNumber,
		Picture:     input.Picture,
	}
}
