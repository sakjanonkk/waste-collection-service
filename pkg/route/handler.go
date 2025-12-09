package route

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

// test
type routeHandler struct {
	service domain.RouteService
}

func NewRouteHandler(router fiber.Router, service domain.RouteService) {
	handler := &routeHandler{service: service}

	router.Post("/", handler.CreateRoute())
	router.Get("/", handler.GetRoutes())
	router.Get("/:id", handler.GetRouteByID())
	// router.Put("/:id", handler.UpdateRoute())
	router.Delete("/:id", handler.DeleteRoute())
}

// CreateRoute godoc
// @Summary Create a new route
// @Description Create a new route with initial points
// @Tags routes
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param route body models.RouteInput true "Route Data"
// @Router /routes [post]
func (h *routeHandler) CreateRoute() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.RouteInput
		if err := c.BodyParser(&input); err != nil {
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

		routeResponse, err := h.service.CreateRoute(input)
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
			Data:    routeResponse,
		})
	}
}

// GetRoutes godoc
// @Summary Get all routes
// @Description Get a list of all routes with pagination and details
// @Tags routes
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Router /routes [get]
func (h *routeHandler) GetRoutes() fiber.Handler {
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

		routes, paginated, err := h.service.GetRoutes(pagination)
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
				"routes":     routes,
				"pagination": paginated,
			},
		})
	}
}

// GetRouteByID godoc
// @Summary Get a route by ID
// @Description Get details of a specific route by ID
// @Tags routes
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Route ID"
// @Router /routes/{id} [get]
func (h *routeHandler) GetRouteByID() fiber.Handler {
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

		routeResponse, err := h.service.GetRouteByID(uint(id))
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
			Data:    routeResponse,
		})
	}
}

// DeleteRoute godoc
// @Summary Delete a route
// @Description Delete a route by ID
// @Tags routes
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "Route ID"
// @Router /routes/{id} [delete]
func (h *routeHandler) DeleteRoute() fiber.Handler {
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

		if err := h.service.DeleteRoute(uint(id)); err != nil {
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
			Data:    "Route deleted successfully",
		})
	}
}
