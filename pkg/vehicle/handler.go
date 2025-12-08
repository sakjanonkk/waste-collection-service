package vehicle

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type vehicleHandler struct {
	service domain.VehicleService
}

func NewVehicleHandler(router fiber.Router, service domain.VehicleService) {
	handler := &vehicleHandler{service: service}

	router.Post("/", handler.CreateVehicle())
	router.Get("/", handler.GetVehicles())
	router.Get("/:id", handler.GetVehicleByID())
	router.Put("/:id", handler.UpdateVehicle())
	router.Delete("/:id", handler.DeleteVehicle())
}

// CreateVehicle godoc
// @Summary Create a new vehicle
// @Description Create a new vehicle with the provided details
// @Tags vehicles
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param vehicle body models.VehicleInput true "Vehicle Data"
// @Success 201 {object} helpers.ResponseForm{data=models.Vehicle}
// @Failure 400 {object} helpers.ResponseForm
// @Failure 409 {object} helpers.ResponseForm
// @Failure 500 {object} helpers.ResponseForm
// @Router /vehicles [post]
func (h *vehicleHandler) CreateVehicle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var vehicleInput models.VehicleInput

		if err := json.Unmarshal(c.Body(), &vehicleInput); err != nil {
			fmt.Printf("❌ Manual JSON Parse Error: %v\n", err)
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "JSON Parse Error",
					Message: err.Error(),
				}},
			})
		}

		vehicleData := vehicleInput.ToVehicle()

		vehicle, err := h.service.CreateVehicle(vehicleData)
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
			Data:    vehicle,
		})
	}
}

// GetVehicles godoc
// @Summary Get all vehicles
// @Description Get a list of all vehicles with pagination
// @Tags vehicles
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} helpers.ResponseForm{data=object{vehicles=[]models.Vehicle,pagination=models.Pagination}}
// @Failure 500 {object} helpers.ResponseForm
// @Router /vehicles [get]
func (h *vehicleHandler) GetVehicles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagination := models.Pagination{
			Page:    1,  // default
			PerPage: 10, // default
		}

		_ = c.QueryParser(&pagination)

		if pagination.PerPage > 100 {
			pagination.PerPage = 100
		}

		vehicles, paginated, err := h.service.GetVehicles(pagination)
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
				"vehicles":   vehicles,
				"pagination": paginated,
			},
		})
	}
}

// GetVehicleByID godoc
// @Summary Get a vehicle by ID
// @Description Get details of a specific vehicle by its ID
// @Tags vehicles
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Vehicle ID"
// @Success 200 {object} helpers.ResponseForm{data=models.Vehicle}
// @Failure 400 {object} helpers.ResponseForm
// @Failure 404 {object} helpers.ResponseForm
// @Router /vehicles/{id} [get]
func (h *vehicleHandler) GetVehicleByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var vehicleInput models.Vehicle
		if err := c.ParamsParser(&vehicleInput); err != nil {
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

		vehicle, err := h.service.GetVehicleByID(vehicleInput)
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
			Data:    vehicle,
		})
	}
}

// UpdateVehicle godoc
// @Summary Update a vehicle
// @Description Update details of an existing vehicle
// @Tags vehicles
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Vehicle ID"
// @Param vehicle body models.VehicleInput true "Vehicle Data"
// @Success 200 {object} helpers.ResponseForm{data=models.Vehicle}
// @Failure 400 {object} helpers.ResponseForm
// @Failure 404 {object} helpers.ResponseForm
// @Failure 409 {object} helpers.ResponseForm
// @Failure 500 {object} helpers.ResponseForm
// @Router /vehicles/{id} [put]
func (h *vehicleHandler) UpdateVehicle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var vehicleInput models.VehicleInput

		if err := json.Unmarshal(c.Body(), &vehicleInput); err != nil {
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

		vehicleData := vehicleInput.ToVehicle()
		vehicleData.ID = uint(parsedID)

		vehicle, err := h.service.UpdateVehicle(vehicleData)
		if err != nil {
			// ✅ แยก error type
			statusCode := fiber.StatusInternalServerError
			title := "Internal Server Error"

			errMsg := err.Error()
			if strings.Contains(errMsg, "not found") {
				statusCode = fiber.StatusNotFound
				title = "Not Found"
			} else if strings.Contains(errMsg, "required") ||
				strings.Contains(errMsg, "invalid") ||
				strings.Contains(errMsg, "must be") ||
				strings.Contains(errMsg, "cannot be negative") {
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

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    vehicle,
		})
	}
}

// DeleteVehicle godoc
// @Summary Delete a vehicle
// @Description Delete a vehicle by its ID
// @Tags vehicles
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Vehicle ID"
// @Success 200 {object} helpers.ResponseForm{data=string}
// @Failure 400 {object} helpers.ResponseForm
// @Failure 404 {object} helpers.ResponseForm
// @Failure 500 {object} helpers.ResponseForm
// @Router /vehicles/{id} [delete]
func (h *vehicleHandler) DeleteVehicle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var vehicleInput models.Vehicle

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
		vehicleInput.ID = uint(parsedID)

		err = h.service.DeleteVehicle(vehicleInput)
		if err != nil {
			// ✅ แยก error type
			statusCode := fiber.StatusInternalServerError
			title := "Internal Server Error"

			if strings.Contains(err.Error(), "not found") {
				statusCode = fiber.StatusNotFound
				title = "Not Found"
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

		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    "Vehicle deleted successfully",
		})
	}
}
