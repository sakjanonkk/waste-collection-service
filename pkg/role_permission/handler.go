package role_permission

import (
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type rolePermissionHandler struct {
	service domain.RolePermissionService
}

func NewRolePermissionHandler(router fiber.Router, resource *handlers.RouterResources, service domain.RolePermissionService) {
	handler := &rolePermissionHandler{service: service}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Create)), handler.CreateRolePermission())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Read)), handler.GetRolePermission())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.List)), handler.GetRolePermissions())
	router.Get("/role/:roleId", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.List)), handler.GetRolePermissionsByRoleID())
	router.Put("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Update)), handler.UpdateRolePermission())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Delete)), handler.DeleteRolePermission())
	router.Delete("/role/:roleId", resource.ReqAuthPerms(models.PermissionGroupName(models.RolePermissionGroup, models.Delete)), handler.DeleteRolePermissionsByRoleID())
}

func (h *rolePermissionHandler) CreateRolePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var rolePermission models.RolePermission
		if err := c.BodyParser(&rolePermission); err != nil {
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
		if errorForm := h.service.CreateRolePermission(rolePermission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermission,
		})
	}
}

func (h *rolePermissionHandler) GetRolePermission() fiber.Handler {
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
		rolePermission, errorForm := h.service.GetRolePermission(uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermission,
		})
	}
}

func (h *rolePermissionHandler) GetRolePermissions() fiber.Handler {
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
		rolePermissions, paginated, searched, errorForm := h.service.GetRolePermissions(pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermissions,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

func (h *rolePermissionHandler) GetRolePermissionsByRoleID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId, err := c.ParamsInt("roleId")
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
		rolePermissions, errorForm := h.service.GetRolePermissionsByRoleID(uint(roleId))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermissions,
		})
	}
}

func (h *rolePermissionHandler) UpdateRolePermission() fiber.Handler {
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
		var rolePermission models.RolePermission
		if err := c.BodyParser(&rolePermission); err != nil {
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
		if errorForm := h.service.UpdateRolePermission(uint(id), rolePermission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    rolePermission,
		})
	}
}

func (h *rolePermissionHandler) DeleteRolePermission() fiber.Handler {
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
		if errorForm := h.service.DeleteRolePermission(uint(id)); errorForm != nil {
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

func (h *rolePermissionHandler) DeleteRolePermissionsByRoleID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId, err := c.ParamsInt("roleId")
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
		if errorForm := h.service.DeleteRolePermissionsByRoleID(uint(roleId)); errorForm != nil {
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
