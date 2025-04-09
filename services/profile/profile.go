package services_profile

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
)

type ProfileService struct {
	DB        *gorm.DB
	AppConfig *bootstrap.AppConfig
}

func NewProfileService(db *gorm.DB, appConfig *bootstrap.AppConfig) *ProfileService {
	return &ProfileService{
		DB:        db,
		AppConfig: appConfig,
	}
}
