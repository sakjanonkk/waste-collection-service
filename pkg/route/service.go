package route

import (
	"errors"

	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type routeService struct {
	repo domain.RouteRepository
}

func NewRouteService(repo domain.RouteRepository) domain.RouteService {
	return &routeService{repo: repo}
}

func (s *routeService) CreateRoute(input models.RouteInput) (models.Route, error) {
	route, err := input.ToRoute()
	if err != nil {
		return models.Route{}, err
	}

	// Validate (basic)
	if route.DriverID == 0 {
		return models.Route{}, errors.New("driver_id is required")
	}
	if route.VehicleID == 0 {
		return models.Route{}, errors.New("vehicle_id is required")
	}

	// Create Route
	createdRoute, err := s.repo.CreateRoute(route)
	if err != nil {
		return models.Route{}, err
	}

	// Create RoutePoints
	for _, pointInput := range input.Points {
		point := models.RoutePoint{
			RouteID:       createdRoute.ID,
			PointID:       pointInput.PointID,
			SequenceOrder: pointInput.SequenceOrder,
			Status:        models.RoutePointStatusPending,
		}
		if _, err := s.repo.AddRoutePoint(point); err != nil {
			return createdRoute, err
		}
	}

	// Return full route
	return s.repo.GetRouteByID(createdRoute.ID)
}

func (s *routeService) GetRoutes(pagination models.Pagination) ([]models.Route, models.Pagination, error) {
	return s.repo.GetRoutes(pagination)
}

func (s *routeService) GetRouteByID(id uint) (models.Route, error) {
	return s.repo.GetRouteByID(id)
}

func (s *routeService) UpdateRoute(route models.Route) (models.Route, error) {
	return s.repo.UpdateRoute(route)
}

func (s *routeService) DeleteRoute(id uint) error {
	return s.repo.DeleteRoute(id)
}
