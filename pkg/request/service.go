package request

import (
	"errors"

	"github.com/zercle/gofiber-skelton/pkg/domain"
	"github.com/zercle/gofiber-skelton/pkg/models"
)

type requestService struct {
	repo domain.RequestRepository
}

func NewRequestService(repo domain.RequestRepository) domain.RequestService {
	return &requestService{repo: repo}
}

func (s *requestService) CreateRequest(input models.RequestInput) (models.RequestResponse, error) {
	req := input.ToRequest()

	// Validate
	if req.RequestType == "" {
		return models.RequestResponse{}, errors.New("request_type is required")
	}
	if req.CreatedByID == nil && req.ReporterName == "" {
		return models.RequestResponse{}, errors.New("created_by_id or reporter_name is required")
	}

	createdReq, err := s.repo.CreateRequest(req)
	if err != nil {
		return models.RequestResponse{}, err
	}

	// Fetch full details to return (relations)
	fullReq, err := s.repo.GetRequestByID(createdReq.ID)
	if err != nil {
		return models.RequestResponse{}, err
	}

	return fullReq.ToResponse(), nil
}

func (s *requestService) GetRequests(pagination models.Pagination, requestType string, status string) ([]models.RequestResponse, models.Pagination, error) {
	filter := make(map[string]interface{})
	if requestType != "" {
		filter["request_type"] = requestType
	}
	if status != "" {
		filter["status"] = status
	}

	requests, paginated, err := s.repo.GetRequests(pagination, filter)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	responses := make([]models.RequestResponse, len(requests))
	for i, r := range requests {
		responses[i] = r.ToResponse()
	}

	return responses, paginated, nil
}

func (s *requestService) GetRequestByID(id uint) (models.RequestResponse, error) {
	req, err := s.repo.GetRequestByID(id)
	if err != nil {
		return models.RequestResponse{}, err
	}
	return req.ToResponse(), nil
}

func (s *requestService) ApproveRequest(id uint) error {
	return s.repo.UpdateRequestStatus(id, models.RequestStatusApproved)
}

func (s *requestService) RejectRequest(id uint) error {
	return s.repo.UpdateRequestStatus(id, models.RequestStatusRejected)
}

func (s *requestService) UpdateRequest(id uint, input models.RequestInput) (models.RequestResponse, error) {
	req, err := s.repo.GetRequestByID(id)
	if err != nil {
		return models.RequestResponse{}, err
	}

	// Update fields
	if input.Remarks != "" {
		req.Remarks = input.Remarks
	}
	// Add other update logic as needed

	updatedReq, err := s.repo.UpdateRequest(req)
	if err != nil {
		return models.RequestResponse{}, err
	}

	return updatedReq.ToResponse(), nil
}

func (s *requestService) DeleteRequest(id uint) error {
	return s.repo.DeleteRequest(id)
}
