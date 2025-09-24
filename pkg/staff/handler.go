package staff

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type staffHandler struct {
	service domain.StaffService
}

func NewStaffHandler(router fiber.Router, service domain.StaffService) {
	handler := &staffHandler{service: service}

	router.Post("/", handler.CreateStaff())
	router.Get("/", handler.GetStaffs())
	router.Get("/:id", handler.GetStaffByID())
	router.Put("/:id", handler.UpdateStaff())    // ✅ update
	router.Delete("/:id", handler.DeleteStaff()) // ✅ delete
}

func (h *staffHandler) CreateStaff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var staffInput models.Staff

		// FE wrong input format
		if err := c.BodyParser(&staffInput); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Bad Request",
						Message: err.Error(),
					},
				},
			})
		}
		staff, err := h.service.CreateStaff(staffInput)

		// BE wrong input format
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusInternalServerError,
						Source:  helpers.WhereAmI(),
						Title:   "Internal Server Error",
						Message: err.Error(),
					},
				},
			})
		}

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    staff,
		})
	}
}

func (h *staffHandler) GetStaffs() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagination models.Pagination
		if err := c.QueryParser(&pagination); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Bad Request",
						Message: err.Error(),
					},
				},
			})
		}

		staffs, paginated, err := h.service.GetStaffs(pagination)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusInternalServerError,
						Source:  helpers.WhereAmI(),
						Title:   "Internal Server Error",
						Message: err.Error(),
					},
				},
			})
		}
		paginated.Page = pagination.Page
		paginated.PerPage = pagination.PerPage

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data: fiber.Map{
				"staffs":     staffs,
				"pagination": paginated,
			},
		})
	}
}

func (h *staffHandler) GetStaffByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var staffInput models.Staff
		if err := c.ParamsParser(&staffInput); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Bad Request",
					Message: err.Error(),
				}},
			})
		}

		staff, err := h.service.GetStaffByID(staffInput)
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

func (h *staffHandler) UpdateStaff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var staffInput models.Staff

		if err := c.BodyParser(&staffInput); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Bad Request",
					Message: err.Error(),
				}},
			})
		}

		parsedID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Invalid ID",
					Message: "ID must be a positive integer.",
				}},
			})
		}
		staffInput.ID = uint(parsedID)

		staff, err := h.service.UpdateStaff(staffInput)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusInternalServerError,
					Source:  helpers.WhereAmI(),
					Title:   "Internal Server Error",
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

func (h *staffHandler) DeleteStaff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var staffInput models.Staff

		parsedID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Invalid ID",
					Message: "ID must be a positive integer.",
				}},
			})
		}
		staffInput.ID = uint(parsedID)

		err = h.service.DeleteStaff(staffInput)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusInternalServerError,
					Source:  helpers.WhereAmI(),
					Title:   "Internal Server Error",
					Message: err.Error(),
				}},
			})
		}

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    "Staff deleted successfully",
		})
	}
}
