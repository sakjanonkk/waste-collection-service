package request

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"github.com/zercle/gofiber-skelton/pkg/utils"
)

type requestHandler struct {
	service domain.RequestService
}

func NewRequestPublicHandler(router fiber.Router, service domain.RequestService) {
	handler := &requestHandler{service: service}
	router.Post("/", handler.CreateRequest())
}

func NewRequestProtectedHandler(router fiber.Router, service domain.RequestService) {
	handler := &requestHandler{service: service}

	// router.Post("/", handler.CreateRequest()) // Public now
	router.Get("/", handler.GetRequests())
	router.Get("/:id", handler.GetRequestByID())
	router.Delete("/:id", handler.DeleteRequest())
	// router.Put("/:id", handler.UpdateRequest()) // Optional
	router.Put("/:id/approve", handler.ApproveRequest())
	router.Put("/:id/reject", handler.RejectRequest())
}

// CreateRequest godoc
// @Summary Create a new request
// @Description Create a new request (report problem or request point) with optional image upload
// @Tags requests
// @Accept  multipart/form-data
// @Produce  json
// @Security ApiKeyAuth
// @Param request_type formData string true "Request Type (report_problem or request_point)"
// @Param point_id formData int false "Point ID"
// @Param point_name formData string false "Point Name"
// @Param latitude formData number false "Latitude"
// @Param longitude formData number false "Longitude"
// @Param remarks formData string false "Remarks"
// @Param created_by_id formData int false "Created By Staff ID"
// @Param reporter_name formData string false "Reporter Name"
// @Param reporter_contact formData string false "Reporter Contact"
// @Param point_image formData file false "Point Image"
// @Router /requests [post]
func (h *requestHandler) CreateRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.RequestInput

		// Parse form fields manually
		input.RequestType = models.RequestType(c.FormValue("request_type"))
		input.PointName = c.FormValue("point_name")
		input.Remarks = c.FormValue("remarks")
		input.ReporterName = c.FormValue("reporter_name")
		input.ReporterContact = c.FormValue("reporter_contact")

		if pointIDStr := c.FormValue("point_id"); pointIDStr != "" {
			if id, err := strconv.ParseUint(pointIDStr, 10, 64); err == nil {
				uid := uint(id)
				input.PointID = &uid
			}
		}

		if latStr := c.FormValue("latitude"); latStr != "" {
			if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
				input.Latitude = lat
			}
		}

		if longStr := c.FormValue("longitude"); longStr != "" {
			if long, err := strconv.ParseFloat(longStr, 64); err == nil {
				input.Longitude = long
			}
		}

		if createdByStr := c.FormValue("created_by_id"); createdByStr != "" {
			if id, err := strconv.ParseUint(createdByStr, 10, 64); err == nil {
				uid := uint(id)
				input.CreatedByID = &uid
			}
		}

		// Handle file upload
		file, err := c.FormFile("point_image")
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
			input.PointImage = url
		}

		resp, err := h.service.CreateRequest(input)
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
			Data:    resp,
		})
	}
}

// GetRequests godoc
// @Summary Get all requests
// @Description Get a list of all requests with pagination and filtering
// @Tags requests
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param request_type query string false "Filter by request type"
// @Param status query string false "Filter by status"
// @Router /requests [get]
func (h *requestHandler) GetRequests() fiber.Handler {
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

		requestType := c.Query("request_type")
		status := c.Query("status")

		requests, paginated, err := h.service.GetRequests(pagination, requestType, status)
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
				"requests":   requests,
				"pagination": paginated,
			},
		})
	}
}

// GetRequestByID godoc
// @Summary Get a request by ID
// @Description Get details of a specific request
// @Tags requests
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Request ID"
// @Router /requests/{id} [get]
func (h *requestHandler) GetRequestByID() fiber.Handler {
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

		req, err := h.service.GetRequestByID(uint(id))
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
			Data:    req,
		})
	}
}

// ApproveRequest godoc
// @Summary Approve a request
// @Description Change request status to approved
// @Tags requests
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Request ID"
// @Router /requests/{id}/approve [put]
func (h *requestHandler) ApproveRequest() fiber.Handler {
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

		if err := h.service.ApproveRequest(uint(id)); err != nil {
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
			Data:    "Request approved successfully",
		})
	}
}

// RejectRequest godoc
// @Summary Reject a request
// @Description Change request status to rejected
// @Tags requests
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Request ID"
// @Router /requests/{id}/reject [put]
func (h *requestHandler) RejectRequest() fiber.Handler {
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

		if err := h.service.RejectRequest(uint(id)); err != nil {
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
			Data:    "Request rejected successfully",
		})
	}
}

// DeleteRequest godoc
// @Summary Delete a request
// @Description Delete a request by ID
// @Tags requests
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Request ID"
// @Router /requests/{id} [delete]
func (h *requestHandler) DeleteRequest() fiber.Handler {
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

		if err := h.service.DeleteRequest(uint(id)); err != nil {
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
			Data:    "Request deleted successfully",
		})
	}
}
