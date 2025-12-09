package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type RouteRepository interface {
	CreateRoute(route models.Route) (models.Route, error)
	GetRoutes(pagination models.Pagination) ([]models.Route, models.Pagination, error)
	GetRouteByID(id uint) (models.Route, error)
	UpdateRoute(route models.Route) (models.Route, error)
	DeleteRoute(id uint) error

	// RoutePoint operations (usually managed via Route, but separate methods can be useful)
	AddRoutePoint(point models.RoutePoint) (models.RoutePoint, error)
	UpdateRoutePointStatus(id uint, status models.RoutePointStatus) error
}

type RouteService interface {
	CreateRoute(input models.RouteInput) (models.RouteResponse, error)
	GetRoutes(pagination models.Pagination) ([]models.RouteResponse, models.Pagination, error)
	GetRouteByID(id uint) (models.RouteResponse, error)
	UpdateRoute(route models.Route) (models.RouteResponse, error)
	DeleteRoute(id uint) error
}
