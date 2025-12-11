package role

import (
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type roleHandler struct {
	service domain.RoleService
}

func NewRoleHandler(router fiber.Router, resource *handlers.RouterResources, service domain.RoleService) {
	handler := &roleHandler{service: service}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Create)), handler.CreateRole())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Read)), handler.GetRole())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.List)), handler.GetRoles())
	router.Put("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Update)), handler.UpdateRole())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RoleGroup, models.Delete)), handler.DeleteRole())
}

func (h *roleHandler) CreateRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var role models.Role
		if err := c.BodyParser(&role); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if errorForm := h.service.CreateRole(role); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    role,
		})
	}
}

func (h *roleHandler) GetRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		role, errorForm := h.service.GetRole(uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    role,
		})
	}
}

func (h *roleHandler) GetRoles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagination := models.Pagination{}
		search := models.Search{}
		if err := c.QueryParser(&pagination); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if err := c.QueryParser(&search); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		roles, paginated, searched, errorForm := h.service.GetRoles(pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    roles,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

func (h *roleHandler) UpdateRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		var role models.Role
		if err := c.BodyParser(&role); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if errorForm := h.service.UpdateRole(uint(id), role); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    role,
		})
	}
}

func (h *roleHandler) DeleteRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		if errorForm := h.service.DeleteRole(uint(id)); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
		})
	}
}
