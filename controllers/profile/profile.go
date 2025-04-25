package controllers_profile

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "کاربر پیدا نشد",
			})
			return
		}
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

func (pc *ProfileController) GetUserAvatarObject(c *gin.Context) {
	fileName := c.Param("file_uuid")

	// Get optional resize height
	heightStr := c.Query("h")
	height, err := strconv.Atoi(heightStr)
	if heightStr != "" {
		if err != nil || height <= 0 || height > 2048 { // Limit max size
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "ارتفاع تصویر نامعتبر است. باید یک عدد بین 1 تا 2048 باشد",
			})
			return
		}
	}

	avatarStream, eTag, err := pc.ProfileService.GetAvatar(c, fileName, height)
	if err != nil {
		switch err {
		case services_profile.ErrInvalidFile:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "فرمت فایل نامعتبر است",
			})
		case services_profile.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "کاربر یا تصویر آواتار پیدا نشد",
			})
		case services_profile.ErrPermissionDenied:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "شما مجاز به مشاهده این تصویر نیستید",
			})
		case services_profile.ErrMinioUnavailable:
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "سرویس ذخیره‌سازی آواتار در دسترس نیست",
			})
			log.Println("Minio service unavailable:", err)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "مشکل غیرمنتظره ای رخ داده است",
			})
			log.Println("Error getting avatar for file", fileName, ":", err)
		}
		return
	}
	defer avatarStream.Close()

	c.Header("ETag", eTag)
	c.Header("Cache-Control", "public, max-age=3600")
	c.Header("Content-Type", "image/png")

	_, err = io.Copy(c.Writer, avatarStream)
	if err != nil {
		log.Println("Error writing avatar using Copy func:", err)
	}
}

func (pc *ProfileController) UploadUserAvatar(c *gin.Context) {
	userID := c.Param("id")
	if _, err := strconv.Atoi(userID); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
			return
	}
	
	resp, err := pc.ProfileService.UploadUserAvatar(c, userID)
	if err != nil {
		if err == services_profile.ErrNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err == services_profile.ErrInvalidSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "سایز فایل نامعتبر است. ماکسیمم مقدار مجاز ۱۶ مگابایت است",
			})
			return
		}
		if err == services_profile.ErrInvalidFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "فرمت فایل نامعتبر است"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusAccepted, resp)
}
