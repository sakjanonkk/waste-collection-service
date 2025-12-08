package staff

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/utils"
)

type staffHandler struct {
	service domain.StaffService
}

func NewStaffHandler(router fiber.Router, service domain.StaffService) {
	handler := &staffHandler{service: service}

	router.Post("/", handler.CreateStaff())
	router.Get("/", handler.GetStaffs())
	router.Get("/:id", handler.GetStaffByID())
	router.Put("/:id", handler.UpdateStaff())
	router.Delete("/:id", handler.DeleteStaff())
}

// CreateStaff godoc
// @Summary Create a new staff member
// @Description Create a new staff member with the provided details. Supports file upload for picture.
// @Tags staff
// @Accept  multipart/form-data
// @Produce  json
// @Security ApiKeyAuth
// @Param prefix formData string false "Prefix"
// @Param firstname formData string true "First Name"
// @Param lastname formData string true "Last Name"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param role formData string true "Role"
// @Param status formData string true "Status"
// @Param phone_number formData string true "Phone Number"
// @Param picture formData file false "Profile Picture"
// @Router /staff [post]
func (h *staffHandler) CreateStaff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var staffInput models.StaffInput

		// Parse form fields manually
		staffInput.Prefix = c.FormValue("prefix")
		staffInput.FirstName = c.FormValue("firstname")
		staffInput.LastName = c.FormValue("lastname")
		staffInput.Email = c.FormValue("email")
		staffInput.Password = c.FormValue("password")
		staffInput.Role = models.StaffRole(c.FormValue("role"))
		staffInput.Status = models.StaffStatus(c.FormValue("status"))
		staffInput.PhoneNumber = c.FormValue("phone_number")

		// Handle file upload
		file, err := c.FormFile("picture")
		if err == nil {
			url, err := utils.UploadFileToMinio(context.Background(), file)
			if err != nil {
				return err
			}
			staffInput.Picture = url
		}

		staff, err := h.service.CreateStaff(staffInput.ToStaff())
		if err != nil {
			statusCode := fiber.StatusInternalServerError
			title := "Internal Server Error"

			errMsg := err.Error()
			if strings.Contains(errMsg, "required") ||
				strings.Contains(errMsg, "invalid") ||
				strings.Contains(errMsg, "must be") {
				statusCode = fiber.StatusBadRequest
				title = "Validation Error"
			} else if strings.Contains(errMsg, "duplicate") ||
				strings.Contains(errMsg, "UNIQUE constraint") {
				statusCode = fiber.StatusConflict
				title = "Duplicate Entry"
			}

			return c.Status(statusCode).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    statusCode,
					Source:  helpers.WhereAmI(),
					Title:   title,
					Message: err.Error(),
				}},
			})
		}

		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    staff,
		})
	}
}

// GetStaffs godoc
// @Summary Get all staff members
// @Description Get a list of all staff members with pagination
// @Tags staff
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Router /staff [get]
func (h *staffHandler) GetStaffs() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var pagination models.Pagination
		if err := c.QueryParser(&pagination); err != nil {
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

		staffs, paginated, err := h.service.GetStaffs(pagination)
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

// GetStaffByID godoc
// @Summary Get a staff member by ID
// @Description Get details of a specific staff member by their ID
// @Tags staff
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Staff ID"
// @Router /staff/{id} [get]
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

// UpdateStaff godoc
// @Summary Update a staff member
// @Description Update details of an existing staff member
// @Tags staff
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Staff ID"
// @Param staff body models.StaffInput true "Staff Data"
// @Router /staff/{id} [put]
func (h *staffHandler) UpdateStaff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var staffInput models.StaffInput

		// ✅ เปลี่ยนเป็น json.Unmarshal()
		if err := json.Unmarshal(c.Body(), &staffInput); err != nil {
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

		staffData := staffInput.ToStaff()
		staffData.ID = uint(parsedID)

		staff, err := h.service.UpdateStaff(staffData)
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

// DeleteStaff godoc
// @Summary Delete a staff member
// @Description Delete a staff member by their ID
// @Tags staff
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Staff ID"
// @Router /staff/{id} [delete]
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
