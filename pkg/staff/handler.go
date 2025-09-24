package staff

import (
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
