package services_pagination

import "gorm.io/gorm"

type PaginationService struct {
	DB *gorm.DB
}

func NewPaginationService(db *gorm.DB) *PaginationService {
	return &PaginationService{
		DB: db,
	}
}
