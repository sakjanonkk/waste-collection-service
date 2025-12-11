package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;index"`
	Group     string    `json:"group" gorm:"column:pkg;not null;index"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
type Role struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;index"`
	Name      string    `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
type RolePermission struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement;index"`
	RoleID       uint       `json:"role_id" gorm:"not null;index"`
	PermissionID uint       `json:"permission_id" gorm:"not null;index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Role         Role       `json:"-" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Permission   Permission `json:"-" gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	RoleID    uint      `json:"role_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Role      Role      `json:"-" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User      Staff     `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func MigrateTablePermissions(db *gorm.DB) error {
	err := db.AutoMigrate(&Permission{}, &Role{}, &RolePermission{}, &UserRole{})
	if err != nil {
		return err
	}
	return err
}

type PermissionGroup string
type PermissionName string

const (
	Create PermissionName = "create"
	Read   PermissionName = "read"
	Update PermissionName = "update"
	Delete PermissionName = "delete"
	List   PermissionName = "list"
)
const (
	UserGroup           PermissionGroup = "user"
	RoleGroup           PermissionGroup = "role"
	PermissionGroup_    PermissionGroup = "permission"
	RolePermissionGroup PermissionGroup = "role_permission"
	UserRoleGroup       PermissionGroup = "user_role"
)

func PermissionGroupName(group PermissionGroup, name PermissionName) string {
	return string(group) + ":" + string(name)
}
