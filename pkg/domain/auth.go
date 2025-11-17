package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type AuthService interface {
	Login(email, password string) (token string, staff models.Staff, err error)
	GetCurrentUser(staffID uint) (staff models.Staff, err error)
	ChangePassword(staffID uint, oldPassword, newPassword string) error
}
