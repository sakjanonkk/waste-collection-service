package utils

import (
	"fmt"
	"strings"

	"github.com/zercle/gofiber-skelton/pkg/models"

	"gorm.io/gorm"
)

func ApplySearch(db *gorm.DB, filter models.Search) *gorm.DB {
	if filter.Keyword == "" || filter.Column == "" {
		return db
	}
	columns := strings.Split(filter.Column, ",")
	var query string

	if len(columns) > 0 {
		var args []interface{}

		for _, column := range columns {
			if query != "" {
				query += " OR "
			}
			query += fmt.Sprintf("%s ILIKE ?", column)
			args = append(args, "%"+filter.Keyword+"%")
		}

		return db.Where(query, args...)
	} else {
		return db.Where(fmt.Sprintf("%s ILIKE ?", filter.Column), "%"+filter.Keyword+"%")
	}
}

func ApplySort(db *gorm.DB, orderBy string, order string) *gorm.DB {
	if orderBy == "" {
		return db
	}
	if order == "" {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderBy, order))
}

func ApplyPagination(db *gorm.DB, pagination *models.Pagination, model interface{}) *gorm.DB {
	var total int64
	db = ApplySort(db, pagination.OrderBy, pagination.Order)
	err := db.Model(model).Count(&total).Error
	if err != nil {
		return nil
	}

	pagination.Total = total
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PerPage < 1 || pagination.PerPage > 1000 {
		pagination.PerPage = 10
	}
	return db.Offset((pagination.Page - 1) * pagination.PerPage).Limit(pagination.PerPage)
}
