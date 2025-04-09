package controllers_profile

import (
	"github.com/gin-gonic/gin"
	services_profile "pashmak.com/pashmak/services/profile"
)

type ProfileController struct {
	ProfileService *services_profile.ProfileService
}

func NewProfileController(profileService *services_profile.ProfileService) *ProfileController{
	return &ProfileController{ProfileService: profileService}
}

func (uc *profileController) GetProfile(c *gin.Context){
	
}
