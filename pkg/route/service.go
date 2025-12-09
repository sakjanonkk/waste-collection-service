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

func (s *routeService) CreateRoute(input models.RouteInput) (models.RouteResponse, error) {
	route, err := input.ToRoute()
	if err != nil {
		return models.RouteResponse{}, err
	}

	// Validate (basic)
	if route.DriverID == 0 {
		return models.RouteResponse{}, errors.New("driver_id is required")
	}
	if route.VehicleID == 0 {
		return models.RouteResponse{}, errors.New("vehicle_id is required")
	}

	// Prepare RoutePoints for atomic creation (GORM handles this if RoutePoints are in the struct)
	if len(input.Points) > 0 {
		route.RoutePoints = make([]models.RoutePoint, len(input.Points))
		for i, p := range input.Points {
			route.RoutePoints[i] = models.RoutePoint{
				PointID:       p.PointID,
				SequenceOrder: p.SequenceOrder,
				Status:        models.RoutePointStatusPending,
			}
		}
	}

	createdRoute, err := s.repo.CreateRoute(route)
	if err != nil {
		return models.RouteResponse{}, err
	}

	// Fetch full details to return
	fullRoute, err := s.repo.GetRouteByID(createdRoute.ID)
	if err != nil {
		return models.RouteResponse{}, err
	}

	return fullRoute.ToResponse(), nil
}

func (s *routeService) GetRoutes(pagination models.Pagination) ([]models.RouteResponse, models.Pagination, error) {
	routes, paginated, err := s.repo.GetRoutes(pagination)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	routeResponses := make([]models.RouteResponse, len(routes))
	for i, r := range routes {
		routeResponses[i] = r.ToResponse()
	}

	return routeResponses, paginated, nil
}

func (s *routeService) GetRouteByID(id uint) (models.RouteResponse, error) {
	route, err := s.repo.GetRouteByID(id)
	if err != nil {
		return models.RouteResponse{}, err
	}
	return route.ToResponse(), nil
}

func (s *routeService) UpdateRoute(route models.Route) (models.RouteResponse, error) {
	updated, err := s.repo.UpdateRoute(route)
	if err != nil {
		return models.RouteResponse{}, err
	}
	return updated.ToResponse(), nil
}

func (s *routeService) DeleteRoute(id uint) error {
	return s.repo.DeleteRoute(id)
}
