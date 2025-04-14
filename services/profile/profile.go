package services_profile

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_auth "pashmak.com/pashmak/models"
	serializers_profile "pashmak.com/pashmak/serializers/profile"
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


func (ps * ProfileService) GetMyProfile(id uint) (serializers_profile.CurrentProfileResponse, error){
	var user models_auth.User
	result := ps.DB.First(&user, "id = ?", id)
	return serializers_profile.CurrentProfileResponse{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Image_url: user.Image_url,
	}, result.Error
}

func (ps *ProfileService) GetProfileByID(id uint)(serializers_profile.GetProfileByIDResponse, error){
	var user models_auth.User
	result := ps.DB.First(&user, "id = ?", id)
	if result.Error != nil{

	}
	return serializers_profile.GetProfileByIDResponse{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Image_url: user.Image_url,
	}, result.Error
}
