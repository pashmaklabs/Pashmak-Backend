package services_profile

import (
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_auth "pashmak.com/pashmak/models/auth"
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
		Avatar_url: user.Avatar_url,
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
		Avatar_url: user.Avatar_url,
	}, result.Error
}

func (ps *ProfileService) GetUserByGmail(email string) (models_auth.User, error){
	var user models_auth.User
	result := ps.DB.First(&user, "email = ?", email)
	return user, result.Error
}

func (ps *ProfileService) UpdateProfile(firstname string, lastname string, avatarUrl string, email string)(error){
	user, err := ps.GetUserByGmail(email)
	if err != nil {
		return err
	}
	user.FirstName = firstname
	user.LastName = lastname
	user.Avatar_url = avatarUrl

	if err := ps.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
