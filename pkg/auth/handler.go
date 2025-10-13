package auth

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type authHandler struct {
	service domain.AuthService
}

func NewAuthHandler(router fiber.Router, service domain.AuthService, jwtResources *models.JwtResources) {
	handler := &authHandler{service: service}

	// Public routes
	router.Post("/login", handler.Login())

	// Protected routes
	router.Get("/me", AuthMiddleware(jwtResources), handler.GetMe())
	router.Put("/change-password", AuthMiddleware(jwtResources), handler.ChangePassword())
}

func (h *authHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loginReq models.LoginRequest

		// Parse request body
		if err := json.Unmarshal(c.Body(), &loginReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Bad Request",
					Message: "Invalid request body: " + err.Error(),
				}},
			})
		}

		// Validate input
		if loginReq.Email == "" || loginReq.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Validation Error",
					Message: "Email and password are required",
				}},
			})
		}

		// Attempt login
		token, staff, err := h.service.Login(loginReq.Email, loginReq.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Title:   "Unauthorized",
					Message: err.Error(),
				}},
			})
		}

		// Return token and staff info
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data: models.LoginResponse{
				Token: token,
				Staff: staff,
			},
		})
	}
}

func (h *authHandler) GetMe() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authenticated user from context
		authUser, ok := GetAuthUser(c)
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

		// Get full staff info
		staff, err := h.service.GetCurrentUser(authUser.StaffID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusNotFound,
					Source:  helpers.WhereAmI(),
					Title:   "Not Found",
					Message: err.Error(),
				}},
			})
		}

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    staff,
		})
	}

}

func (h *authHandler) ChangePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authUser, _ := GetAuthUser(c)

		var req struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		// ✅ เปลี่ยนเป็น json.Unmarshal
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Bad Request",
					Message: "Invalid request body: " + err.Error(),
				}},
			})
		}

		// Validate
		if req.OldPassword == "" || req.NewPassword == "" {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Validation Error",
					Message: "Old password and new password are required",
				}},
			})
		}

		if len(req.NewPassword) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Validation Error",
					Message: "New password must be at least 8 characters",
				}},
			})
		}

		if err := h.service.ChangePassword(authUser.StaffID, req.OldPassword, req.NewPassword); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Change Password Failed",
					Message: err.Error(),
				}},
			})
		}

		return c.JSON(helpers.ResponseForm{
			Success: true,
			Data:    "Password changed successfully",
		})
	}
}
