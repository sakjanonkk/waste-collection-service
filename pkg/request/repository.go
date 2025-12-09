package request

import (
	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
	"gorm.io/gorm"
)

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) domain.RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) CreateRequest(req models.Request) (models.Request, error) {
	err := r.db.Create(&req).Error
	return req, err
}

func (r *requestRepository) GetRequests(pagination models.Pagination, filter map[string]interface{}) (requests []models.Request, paginated models.Pagination, err error) {
	var total int64
	query := r.db.Model(&models.Request{})

	// Apply filters
	if requestType, ok := filter["request_type"]; ok && requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}
	if status, ok := filter["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	err = query.Count(&total).Error
	if err != nil {
		return nil, models.Pagination{}, err
	}
	paginated.Total = total

	err = query.Preload("CreatedBy").Preload("Point").Preload("ApprovedPoint").
		Order("created_at desc").
		Limit(pagination.PerPage).Offset((pagination.Page - 1) * pagination.PerPage).
		Find(&requests).Error

	if err != nil {
		return nil, models.Pagination{}, err
	}
	return requests, paginated, nil
}

func (r *requestRepository) GetRequestByID(id uint) (req models.Request, err error) {
	err = r.db.Preload("CreatedBy").Preload("Point").Preload("ApprovedPoint").First(&req, id).Error
	return req, err
}

func (r *requestRepository) UpdateRequestStatus(id uint, status models.RequestStatus) error {
	return r.db.Model(&models.Request{}).Where("id = ?", id).Update("status", status).Error
}

func (r *requestRepository) UpdateRequest(req models.Request) (models.Request, error) {
	err := r.db.Model(&models.Request{}).Where("id = ?", req.ID).Updates(req).Error
	if err != nil {
		return req, err
	}
	return r.GetRequestByID(req.ID)
}

func (r *requestRepository) DeleteRequest(id uint) error {
	return r.db.Delete(&models.Request{}, id).Error
}
