package permission

import (
	"github.com/zercle/gofiber-skelton/internal/handlers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type permissionHandler struct {
	service domain.PermissionService
}

func NewPermissionHandler(router fiber.Router, resource *handlers.RouterResources, service domain.PermissionService) {
	handler := &permissionHandler{service: service}
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup_, models.Create)), handler.CreatePermission())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup_, models.Read)), handler.GetPermission())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup_, models.List)), handler.GetPermissions())
	router.Put("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup_, models.Update)), handler.UpdatePermission())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.PermissionGroup_, models.Delete)), handler.DeletePermission())
}

func (h *permissionHandler) CreatePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var permission models.Permission
		if err := c.BodyParser(&permission); err != nil {
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
		if errorForm := h.service.CreatePermission(permission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permission,
		})
	}
}

func (h *permissionHandler) GetPermission() fiber.Handler {
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
		permission, errorForm := h.service.GetPermission(uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permission,
		})
	}
}

func (h *permissionHandler) GetPermissions() fiber.Handler {
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
		permissions, paginated, searched, errorForm := h.service.GetPermissions(pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permissions,
			Result: map[string]interface{}{
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

func (h *permissionHandler) UpdatePermission() fiber.Handler {
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
		var permission models.Permission
		if err := c.BodyParser(&permission); err != nil {
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
		if errorForm := h.service.UpdatePermission(uint(id), permission); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    permission,
		})
	}
}

func (h *permissionHandler) DeletePermission() fiber.Handler {
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
		if errorForm := h.service.DeletePermission(uint(id)); errorForm != nil {
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
