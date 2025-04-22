package controllers_profile

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	serializers_profile "pashmak.com/pashmak/serializers/profile"
	services_auth "pashmak.com/pashmak/services/auth"
	services_profile "pashmak.com/pashmak/services/profile"
)

type ProfileController struct {
	ProfileService *services_profile.ProfileService
}

func NewProfileController(profileService *services_profile.ProfileService) *ProfileController {
	return &ProfileController{ProfileService: profileService}
}

func (pc *ProfileController) GetMyProfile(c *gin.Context) {
	value, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	userinfo := value.(services_auth.UserInfo)
	profile, err := pc.ProfileService.GetMyProfile(userinfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, profile)
}

func (pc *ProfileController) GetProfileByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	profile, err := pc.ProfileService.GetProfileByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": profile,
	})
}

func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	var body serializers_profile.UpdateProfileRequest

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}
	// userinfo, exists := c.Get("user")

	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status":  "error",
	// 		"message": "شمامجاز به انجام این عملیات نمی باشید.",
	// 	})
	// 	return
	// }
	// userpayload := userinfo.(services_auth.UserInfo)
	err := pc.ProfileService.UpdateProfile(body.FirstName, body.LastName, body.Avatar_url, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیر منتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "آپدیت با موفقیت انجام شد.",
	})
}
