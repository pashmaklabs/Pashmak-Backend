package controllers_profile

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	services_auth "pashmak.com/pashmak/services/auth"
	services_profile "pashmak.com/pashmak/services/profile"
)

type ProfileController struct {
	ProfileService *services_profile.ProfileService
}

func NewProfileController(profileService *services_profile.ProfileService) *ProfileController {
	return &ProfileController{ProfileService: profileService}
}

func (pc *ProfileController) GetProfile(c *gin.Context) {
	value, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
	}
	userinfo := value.(services_auth.UserInfo)
	profile, err := pc.ProfileService.GetProfile(userinfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
	}
	c.IndentedJSON(http.StatusOK, profile)
}
