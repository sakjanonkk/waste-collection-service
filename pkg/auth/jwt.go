package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

// JwtClaims custom claims for JWT
type JwtClaims struct {
	StaffID uint               `json:"staff_id"`
	Email   string             `json:"email"`
	Role    models.StaffRole   `json:"role"`
	Status  models.StaffStatus `json:"status"`
	jwt.RegisteredClaims
}

// GenerateToken generates JWT token for staff
func GenerateToken(jwtResources *models.JwtResources, staff models.Staff) (string, error) {
	claims := JwtClaims{
		StaffID: staff.ID,
		Email:   staff.Email,
		Role:    staff.Role,
		Status:  staff.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "waste-management-api",
		},
	}

	token := jwt.NewWithClaims(jwtResources.JwtSigningMethod, claims)
	return token.SignedString(jwtResources.JwtSignKey)
}

// ValidateToken validates JWT token and returns claims
func ValidateToken(jwtResources *models.JwtResources, tokenString string) (*JwtClaims, error) {
	token, err := jwtResources.JwtParser.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		jwtResources.JwtKeyfunc,
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
