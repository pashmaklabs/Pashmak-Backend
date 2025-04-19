package services_navigation

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
)

type NavigationService struct {
	DB        *gorm.DB
	AppConfig *bootstrap.AppConfig
}

func NewNavigationService(db *gorm.DB, appconfig *bootstrap.AppConfig) *NavigationService {
	return &NavigationService{
		DB:        db,
		AppConfig: appconfig,
	}
}
