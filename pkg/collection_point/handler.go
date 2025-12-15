package collection_point

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/utils"
)

// test
type collectionPointHandler struct {
	service domain.CollectionPointService
}

func NewCollectionPointHandler(router fiber.Router, service domain.CollectionPointService) {
	handler := &collectionPointHandler{service: service}

	router.Post("/", handler.CreateCollectionPoint())
	router.Get("/", handler.GetCollectionPoints())
	router.Get("/:id", handler.GetCollectionPointByID())
	router.Put("/:id", handler.UpdateCollectionPoint())
	router.Delete("/:id", handler.DeleteCollectionPoint())
}

// CreateCollectionPoint godoc
// @Summary Create a new collection point
// @Description Create a new collection point with optional image
// @Tags collection-points
// @Accept  multipart/form-data
// @Produce  json
// @Security ApiKeyAuth
// @Param name formData string true "Name"
// @Param latitude formData number true "Latitude"
// @Param longitude formData number true "Longitude"
// @Param address formData string false "Address"
// @Param status formData string true "Status (active, inactive)"
// @Param problem_reported formData string false "Problem Reported"
// @Param regular_capacity formData number false "Regular Capacity"
// @Param recycle_capacity formData number false "Recycle Capacity"
// @Param image formData file false "Image"
// @Router /collection-points [post]
func (h *collectionPointHandler) CreateCollectionPoint() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.CollectionPointInput

		// Parse form fields
		input.Name = c.FormValue("name")
		input.Address = c.FormValue("address")
		input.Status = models.CollectionPointStatus(c.FormValue("status"))
		input.ProblemReported = c.FormValue("problem_reported")

		if val := c.FormValue("latitude"); val != "" {
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				input.Latitude = f
			}
		}

		if val := c.FormValue("longitude"); val != "" {
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				input.Longitude = f
			}
		}

		if val := c.FormValue("regular_capacity"); val != "" {
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				input.RegularCapacity = f
			}
		}

		if val := c.FormValue("recycle_capacity"); val != "" {
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				input.RecycleCapacity = f
			}
		}

		// Handle file upload
		file, err := c.FormFile("image")
		if err == nil {
			url, err := utils.UploadFileToMinio(context.Background(), file)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseForm{
					Success: false,
					Errors: []helpers.ResponseError{{
						Code:    fiber.StatusInternalServerError,
						Source:  helpers.WhereAmI(),
						Title:   "Upload Error",
						Message: "Failed to upload image: " + err.Error(),
					}},
				})
			}
			input.Image = url
		}

		point, err := h.service.CreateCollectionPoint(input.ToCollectionPoint())
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

		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
			Data:    point,
		})
	}
}

// GetCollectionPoints godoc
// @Summary Get all collection points
// @Description Get a list of all collection points with pagination
// @Tags collection-points
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Router /collection-points [get]
func (h *collectionPointHandler) GetCollectionPoints() fiber.Handler {
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

		points, paginated, err := h.service.GetCollectionPoints(pagination)
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
			Data: fiber.Map{
				"collection_points": points,
				"pagination":        paginated,
			},
		})
	}
}

// GetCollectionPointByID godoc
// @Summary Get a collection point by ID
// @Description Get details of a specific collection point by ID
// @Tags collection-points
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Collection Point ID"
// @Router /collection-points/{id} [get]
func (h *collectionPointHandler) GetCollectionPointByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Invalid ID",
					Message: "ID must be a positive integer",
				}},
			})
		}

		point, err := h.service.GetCollectionPointByID(uint(id))
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
			Data:    point,
		})
	}
}

// UpdateCollectionPoint godoc
// @Summary Update a collection point
// @Description Update details of an existing collection point
// @Tags collection-points
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Collection Point ID"
// @Param collection_point body models.CollectionPointInput true "Collection Point Data"
// @Router /collection-points/{id} [put]
func (h *collectionPointHandler) UpdateCollectionPoint() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Invalid ID",
					Message: "ID must be a positive integer",
				}},
			})
		}

		var input models.CollectionPointInput
		if err := json.Unmarshal(c.Body(), &input); err != nil {
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

		point := input.ToCollectionPoint()
		point.ID = uint(id)

		updatedPoint, err := h.service.UpdateCollectionPoint(point)
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
			Data:    updatedPoint,
		})
	}
}

// DeleteCollectionPoint godoc
// @Summary Delete a collection point
// @Description Delete a collection point by ID
// @Tags collection-points
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Collection Point ID"
// @Router /collection-points/{id} [delete]
func (h *collectionPointHandler) DeleteCollectionPoint() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{{
					Code:    fiber.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Title:   "Invalid ID",
					Message: "ID must be a positive integer",
				}},
			})
		}

		if err := h.service.DeleteCollectionPoint(uint(id)); err != nil {
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
			Data:    "Collection Point deleted successfully",
		})
	}
}
