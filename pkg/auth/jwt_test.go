package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Setup
	jwtResources := &models.JwtResources{
		JwtSigningMethod: jwt.SigningMethodHS256,
		JwtSignKey:       []byte("secret"),
		JwtKeyfunc: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		JwtParser: jwt.NewParser(),
	}

	staff := models.Staff{
		ID:     1,
		Email:  "test@example.com",
		Role:   models.RoleAdmin,
		Status: models.StatusStaffActive,
	}

	// Test GenerateToken
	tokenString, err := GenerateToken(jwtResources, staff)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Test ValidateToken
	claims, err := ValidateToken(jwtResources, tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Assert Claims (This is where it should fail before the fix)
	assert.Equal(t, staff.ID, claims.StaffID, "StaffID mismatch")
	assert.Equal(t, staff.Email, claims.Email, "Email mismatch")
	assert.Equal(t, staff.Role, claims.Role, "Role mismatch")
	assert.Equal(t, staff.Status, claims.Status, "Status mismatch")
}
