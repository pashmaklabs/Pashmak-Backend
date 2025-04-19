package controllers_navigation

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	services_navigation "pashmak.com/pashmak/services/navigation"
)

type NavigationController struct {
	NavigationService *services_navigation.NavigationService
	DB                *gorm.DB
	AppConfig         *bootstrap.AppConfig
}

func NewNavigationController(navigationService *services_navigation.NavigationService, db *gorm.DB, appConfig *bootstrap.AppConfig) *NavigationController {
	return &NavigationController{
		NavigationService: navigationService,
		DB:                db,
		AppConfig:         appConfig,
	}
}
