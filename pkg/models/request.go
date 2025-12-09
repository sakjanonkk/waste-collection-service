package models

import (
	"time"

	"gorm.io/gorm"
)

type RequestType string
type RequestStatus string

const (
	RequestTypeReportProblem RequestType = "report_problem"
	RequestTypeRequestPoint  RequestType = "request_point"

	RequestStatusPending  RequestStatus = "pending"
	RequestStatusApproved RequestStatus = "approved"
	RequestStatusRejected RequestStatus = "rejected"
)

type Request struct {
	ID              uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	RequestType     RequestType      `json:"request_type" gorm:"column:request_type;not null"`
	PointID         *uint            `json:"point_id,omitempty" gorm:"column:point_id"`
	Point           *CollectionPoint `json:"point,omitempty" gorm:"foreignKey:PointID"`
	PointName       string           `json:"point_name" gorm:"column:point_name"`
	PointImage      string           `json:"point_image" gorm:"column:point_image"`
	Latitude        float64          `json:"latitude" gorm:"column:latitude"`
	Longitude       float64          `json:"longitude" gorm:"column:longitude"`
	Remarks         string           `json:"remarks" gorm:"column:remarks"`
	RequestDatetime time.Time        `json:"request_datetime" gorm:"column:request_datetime;not null"`
	Status          RequestStatus    `json:"status" gorm:"column:status;not null;default:'pending'"`
	CreatedByID     *uint            `json:"created_by_id" gorm:"column:created_by_id"`
	CreatedBy       *Staff           `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	ReporterName    string           `json:"reporter_name" gorm:"column:reporter_name"`
	ReporterContact string           `json:"reporter_contact" gorm:"column:reporter_contact"`
	ApprovedPointID *uint            `json:"approved_point_id,omitempty" gorm:"column:approved_point_id"`
	ApprovedPoint   *CollectionPoint `json:"approved_point,omitempty" gorm:"foreignKey:ApprovedPointID"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt   `json:"-" gorm:"index"`
}

type RequestInput struct {
	RequestType     RequestType `json:"request_type" example:"report_problem"`
	PointID         *uint       `json:"point_id,omitempty" example:"1"`
	PointName       string      `json:"point_name,omitempty" example:"Near Central Park"`
	PointImage      string      `json:"point_image,omitempty" example:"https://example.com/image.jpg"`
	Latitude        float64     `json:"latitude,omitempty" example:"13.7563"`
	Longitude       float64     `json:"longitude,omitempty" example:"100.5018"`
	Remarks         string      `json:"remarks,omitempty" example:"Bin is full and overflowing"`
	CreatedByID     *uint       `json:"created_by_id,omitempty" example:"1"`
	ReporterName    string      `json:"reporter_name,omitempty" example:"John Doe"`
	ReporterContact string      `json:"reporter_contact,omitempty" example:"0812345678"`
}

func (input *RequestInput) ToRequest() Request {
	return Request{
		RequestType:     input.RequestType,
		PointID:         input.PointID,
		PointName:       input.PointName,
		PointImage:      input.PointImage,
		Latitude:        input.Latitude,
		Longitude:       input.Longitude,
		Remarks:         input.Remarks,
		RequestDatetime: time.Now(),
		CreatedByID:     input.CreatedByID,
		ReporterName:    input.ReporterName,
		ReporterContact: input.ReporterContact,
		Status:          RequestStatusPending,
	}
}

type RequestResponse struct {
	ID              uint          `json:"id"`
	RequestType     RequestType   `json:"request_type"`
	RequestTypeName string        `json:"request_type_name"` // e.g. "Report Problem"
	PointID         *uint         `json:"point_id,omitempty"`
	PointName       string        `json:"point_name"`
	PointImage      string        `json:"point_image"`
	Latitude        float64       `json:"latitude"`
	Longitude       float64       `json:"longitude"`
	Remarks         string        `json:"remarks"`
	RequestDatetime string        `json:"request_datetime"`
	Status          RequestStatus `json:"status"`
	CreatedBy       string        `json:"created_by,omitempty"`
	ReporterName    string        `json:"reporter_name,omitempty"`
	ApprovedPointID *uint         `json:"approved_point_id,omitempty"`
}

func (r *Request) ToResponse() RequestResponse {
	resp := RequestResponse{
		ID:              r.ID,
		RequestType:     r.RequestType,
		PointID:         r.PointID,
		PointName:       r.PointName,
		PointImage:      r.PointImage,
		Latitude:        r.Latitude,
		Longitude:       r.Longitude,
		Remarks:         r.Remarks,
		RequestDatetime: r.RequestDatetime.Format(time.RFC3339),
		Status:          r.Status,
		ApprovedPointID: r.ApprovedPointID,
		ReporterName:    r.ReporterName,
	}

	if r.RequestType == RequestTypeReportProblem {
		resp.RequestTypeName = "Report Problem"
	} else if r.RequestType == RequestTypeRequestPoint {
		resp.RequestTypeName = "Request New Point"
	}

	if r.CreatedBy != nil {
		resp.CreatedBy = r.CreatedBy.FirstName + " " + r.CreatedBy.LastName
	}

	// Auto-fill point name from relations if empty in request itself
	if r.PointName == "" && r.Point != nil {
		resp.PointName = r.Point.Name
	}

	return resp
}
