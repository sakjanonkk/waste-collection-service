package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type CollectionPointRepository interface {
	CreateCollectionPoint(collectionPoint models.CollectionPoint) (models.CollectionPoint, error)
	GetCollectionPoints(pagination models.Pagination) ([]models.CollectionPoint, models.Pagination, error)
	GetCollectionPointByID(id uint) (models.CollectionPoint, error)
	UpdateCollectionPoint(collectionPoint models.CollectionPoint) (models.CollectionPoint, error)
	DeleteCollectionPoint(id uint) error
}

type CollectionPointService interface {
	CreateCollectionPoint(collectionPoint models.CollectionPoint) (models.CollectionPoint, error)
	GetCollectionPoints(pagination models.Pagination) ([]models.CollectionPoint, models.Pagination, error)
	GetCollectionPointByID(id uint) (models.CollectionPoint, error)
	UpdateCollectionPoint(collectionPoint models.CollectionPoint) (models.CollectionPoint, error)
	DeleteCollectionPoint(id uint) error
}
