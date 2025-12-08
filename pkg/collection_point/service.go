package collection_point

import (
	"errors"
	"strings"

	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type collectionPointService struct {
	repo domain.CollectionPointRepository
}

func NewCollectionPointService(repo domain.CollectionPointRepository) domain.CollectionPointService {
	return &collectionPointService{repo: repo}
}

func (s *collectionPointService) CreateCollectionPoint(point models.CollectionPoint) (models.CollectionPoint, error) {
	if err := s.validateInput(point); err != nil {
		return point, err
	}
	return s.repo.CreateCollectionPoint(point)
}

func (s *collectionPointService) validateInput(point models.CollectionPoint) error {
	if strings.TrimSpace(point.Name) == "" {
		return errors.New("name is required")
	}
	if point.Latitude == 0 || point.Longitude == 0 {
		return errors.New("valid latitude and longitude are required")
	}
	return nil
}

func (s *collectionPointService) GetCollectionPoints(pagination models.Pagination) ([]models.CollectionPoint, models.Pagination, error) {
	return s.repo.GetCollectionPoints(pagination)
}

func (s *collectionPointService) GetCollectionPointByID(id uint) (models.CollectionPoint, error) {
	return s.repo.GetCollectionPointByID(id)
}

func (s *collectionPointService) UpdateCollectionPoint(point models.CollectionPoint) (models.CollectionPoint, error) {
	return s.repo.UpdateCollectionPoint(point)
}

func (s *collectionPointService) DeleteCollectionPoint(id uint) error {
	return s.repo.DeleteCollectionPoint(id)
}
