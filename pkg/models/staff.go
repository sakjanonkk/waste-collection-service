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

// ✅ Main Staff model (สำหรับ DB และ Response)
type Staff struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix      string         `json:"prefix" gorm:"column:prefix;not null"`
	FirstName   string         `json:"firstname" gorm:"column:firstname;not null"`
	LastName    string         `json:"lastname" gorm:"column:lastname;not null"`
	Email       string         `json:"email" gorm:"column:email;not null;unique"`
	Password    string         `json:"-" gorm:"column:password;not null"` // ไม่ return ใน response
	Role        StaffRole      `json:"role" gorm:"column:role;not null"`
	Status      StaffStatus    `json:"status" gorm:"column:status;not null"`
	PhoneNumber string         `json:"phone_number" gorm:"column:phone;not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ✅ Input struct สำหรับ Create/Update (รับ password ได้)
type StaffInput struct {
	Prefix      string      `json:"prefix"`
	FirstName   string      `json:"firstname"`
	LastName    string      `json:"lastname"`
	Email       string      `json:"email"`
	Password    string      `json:"password"` // ✅ รับ password ได้
	Role        StaffRole   `json:"role"`
	Status      StaffStatus `json:"status"`
	PhoneNumber string      `json:"phone_number"`
}

// Helper: แปลง StaffInput → Staff
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
	}
}
