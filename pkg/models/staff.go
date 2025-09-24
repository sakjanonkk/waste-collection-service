package models

import (
	"time"

	"gorm.io/gorm"
)

type Staff struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix    string         `json:"prefix" gorm:"column:s_prefix;not null"`
	FirstName string         `json:"firstname" gorm:"column:s_firstname;not null"`
	LastName  string         `json:"lastname" gorm:"column:s_lastname;not null"`
	Email     string         `json:"email" gorm:"column:s_email;not null;unique"`
	Password  string         `json:"password" gorm:"column:s_password;not null"`
	Role      string         `json:"role" gorm:"column:s_role;not null"`
	Status    string         `json:"status" gorm:"column:s_status;not null"`
	Phone     string         `json:"phone_number" gorm:"column:s_phone;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
