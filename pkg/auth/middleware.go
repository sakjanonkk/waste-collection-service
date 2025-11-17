package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(jwtResources *models.JwtResources) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Title:   "Unauthorized",
					Message: "Missing authorization token",
				}},
			})
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Title:   "Unauthorized",
					Message: "Invalid authorization format. Use: Bearer <token>",
				}},
			})
		}

		tokenString := parts[1]

		claims, err := ValidateToken(jwtResources, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Title:   "Unauthorized",
					Message: "Invalid or expired token",
				}},
			})
		}

		// Check if staff is active
		if claims.Status != models.StatusStaffActive {
			return c.Status(fiber.StatusForbidden).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusForbidden,
					Source:  helpers.WhereAmI(),
					Title:   "Forbidden",
					Message: "Account is not active",
				}},
			})
		}

		// Store user info in context
		c.Locals("auth_user", models.AuthUser{
			StaffID: claims.StaffID,
			Email:   claims.Email,
			Role:    claims.Role,
			Status:  claims.Status,
		})

		return c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(allowedRoles ...models.StaffRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authenticated user from context
		authUser, ok := c.Locals("auth_user").(models.AuthUser)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Title:   "Unauthorized",
					Message: "Authentication required",
				}},
			})
		}

		// Check if user has allowed role
		hasRole := false
		for _, role := range allowedRoles {
			if authUser.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusForbidden,
					Source:  helpers.WhereAmI(),
					Title:   "Forbidden",
					Message: "You don't have permission to access this resource",
				}},
			})
		}

		return c.Next()
	}
}

// GetAuthUser helper to get authenticated user from context
func GetAuthUser(c *fiber.Ctx) (models.AuthUser, bool) {
	authUser, ok := c.Locals("auth_user").(models.AuthUser)
	return authUser, ok
}

// RequireOwnerOrRole middleware checks if user is owner or has required role
func RequireOwnerOrRole(allowedRoles ...models.StaffRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authUser, ok := c.Locals("auth_user").(models.AuthUser)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Title:   "Unauthorized",
					Message: "Authentication required",
				}},
			})
		}

		// Check if user has allowed role
		for _, role := range allowedRoles {
			if authUser.Role == role {
				return c.Next() // Has permission
			}
		}

		// Check if user is the owner (updating own profile)
		id := c.Params("id")
		if id != "" {
			var staffID uint
			_, err := fmt.Sscanf(id, "%d", &staffID)
			if err == nil && staffID == authUser.StaffID {
				return c.Next() // Is owner
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(helpers.ResponseForm{
			Success: false,
			Errors: []helpers.ResponseError{{
				Code:    fiber.StatusForbidden,
				Source:  helpers.WhereAmI(),
				Title:   "Forbidden",
				Message: "You don't have permission to access this resource",
			}},
		})
	}
}
