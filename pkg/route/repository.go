package route

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type routeRepository struct {
	db *gorm.DB
}

func NewRouteRepository(db *gorm.DB) domain.RouteRepository {
	return &routeRepository{db: db}
}

func (r *routeRepository) CreateRoute(route models.Route) (models.Route, error) {
	err := r.db.Create(&route).Error
	return route, err
}

func (r *routeRepository) GetRoutes(pagination models.Pagination) (routes []models.Route, paginated models.Pagination, err error) {
	var total int64
	err = r.db.Model(&models.Route{}).Count(&total).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	paginated.Total = total

	// Preload driver, vehicle, and points
	err = r.db.Preload("Driver").Preload("Vehicle").Preload("RoutePoints.Point").
		Order("created_at desc").
		Limit(pagination.PerPage).Offset((pagination.Page - 1) * pagination.PerPage).
		Find(&routes).Error

	if err != nil {
		return nil, models.Pagination{}, err
	}
	return routes, paginated, nil
}

func (r *routeRepository) GetRouteByID(id uint) (route models.Route, err error) {
	err = r.db.Preload("Driver").Preload("Vehicle").Preload("RoutePoints.Point").First(&route, id).Error
	return route, err
}

func (r *routeRepository) UpdateRoute(route models.Route) (models.Route, error) {
	err := r.db.Model(&models.Route{}).Where("id = ?", route.ID).Updates(route).Error
	if err != nil {
		return route, err
	}
	return r.GetRouteByID(route.ID)
}

func (r *routeRepository) DeleteRoute(id uint) error {
	return r.db.Delete(&models.Route{}, id).Error
}

func (r *routeRepository) AddRoutePoint(point models.RoutePoint) (models.RoutePoint, error) {
	err := r.db.Create(&point).Error
	return point, err
}

func (r *routeRepository) UpdateRoutePointStatus(id uint, status models.RoutePointStatus) error {
	return r.db.Model(&models.RoutePoint{}).Where("id = ?", id).Update("status", status).Error
}
