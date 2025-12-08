package collection_point

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type collectionPointRepository struct {
	db *gorm.DB
}

func NewCollectionPointRepository(db *gorm.DB) domain.CollectionPointRepository {
	return &collectionPointRepository{db: db}
}

func (r *collectionPointRepository) CreateCollectionPoint(point models.CollectionPoint) (models.CollectionPoint, error) {
	err := r.db.Create(&point).Error
	return point, err
}

func (r *collectionPointRepository) GetCollectionPoints(pagination models.Pagination) (points []models.CollectionPoint, paginated models.Pagination, err error) {
	var total int64
	err = r.db.Model(&models.CollectionPoint{}).Count(&total).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	paginated.Total = total
	err = r.db.Limit(pagination.PerPage).Offset((pagination.Page - 1) * pagination.PerPage).Find(&points).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	return points, paginated, nil
}

func (r *collectionPointRepository) GetCollectionPointByID(id uint) (point models.CollectionPoint, err error) {
	err = r.db.First(&point, id).Error
	return point, err
}

func (r *collectionPointRepository) UpdateCollectionPoint(point models.CollectionPoint) (models.CollectionPoint, error) {
	err := r.db.Model(&models.CollectionPoint{}).Where("id = ?", point.ID).Updates(point).Error
	if err != nil {
		return point, err
	}
	err = r.db.First(&point, point.ID).Error
	return point, err
}

func (r *collectionPointRepository) DeleteCollectionPoint(id uint) error {
	return r.db.Delete(&models.CollectionPoint{}, id).Error
}
