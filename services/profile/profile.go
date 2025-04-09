package services_profile

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_auth "pashmak.com/pashmak/models"
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


func (ps * ProfileService) GetProfile(id uint) (models_auth.User, error){
	var user models_auth.User
	result := ps.DB.First(&user, "id = ?", id)
	return user, result.Error
}
