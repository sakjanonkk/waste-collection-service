package domain

import "github.com/zercle/gofiber-skelton/pkg/models"

type RequestRepository interface {
	CreateRequest(req models.Request) (models.Request, error)
	GetRequests(pagination models.Pagination, filter map[string]interface{}) ([]models.Request, models.Pagination, error)
	GetRequestByID(id uint) (models.Request, error)
	UpdateRequestStatus(id uint, status models.RequestStatus) error
	UpdateRequest(req models.Request) (models.Request, error)
	DeleteRequest(id uint) error
}

type RequestService interface {
	CreateRequest(input models.RequestInput) (models.RequestResponse, error)
	GetRequests(pagination models.Pagination, requestType string, status string) ([]models.RequestResponse, models.Pagination, error)
	GetRequestByID(id uint) (models.RequestResponse, error)
	ApproveRequest(id uint) error
	RejectRequest(id uint) error
	UpdateRequest(id uint, input models.RequestInput) (models.RequestResponse, error)
	DeleteRequest(id uint) error
}
