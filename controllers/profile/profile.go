package controllers_profile

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	serializers_profile "pashmak.com/pashmak/serializers/profile"
	services_auth "pashmak.com/pashmak/services/auth"
	services_place "pashmak.com/pashmak/services/place"
	services_profile "pashmak.com/pashmak/services/profile"
)

type ProfileController struct {
	ProfileService *services_profile.ProfileService
	PlaceService *services_place.PlaceService
}

func NewProfileController(profileService *services_profile.ProfileService, placeService *services_place.PlaceService) *ProfileController {
	return &ProfileController{
		ProfileService: profileService,
		PlaceService: placeService,
	}
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
			c.Status(http.StatusNotFound)
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
	c.Header("Content-Type", "image/webp")

	_, err = io.Copy(c.Writer, avatarStream)
	if err != nil {
		log.Println("Error writing avatar using Copy func:", err)
	}
}

func (pc *ProfileController) UploadUserAvatar(c *gin.Context) {
	value, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	userinfo := value.(services_auth.UserInfo)
	// userID := c.Param("id")
	// if _, err := strconv.Atoi(userID); err != nil {
	// 	c.AbortWithStatus(http.StatusNotFound)
	// 		return
	// }
	
	resp, err := pc.ProfileService.UploadUserAvatar(c, userinfo.ID)
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
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "فرمت فایل نامعتبر است. لطفا یک فایل تصویری ارسال کنید",
			})
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

func (pc *ProfileController) UpdateUserProfile(c *gin.Context){
	validatedData, exists := c.Get("validated")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed to retrieve validated data from context: exists=%v", exists)
		return
	}

	body, ok := validatedData.(serializers_profile.UpdateUserProfileRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed type assertion for validated data: expected UpdateUserProfileRequest, got %T", validatedData)
		return
	}
	userinfo, exists := c.Get("user")
	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)
	err := pc.ProfileService.UpdateUserProfile(userpayload, body)
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
		"message": "آپدیت پروفایل با موفقیت انجام شد",
	})
}

func (pc *ProfileController) FetchSearchHistory(c *gin.Context){
	userinfo, exists := c.Get("user")
	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)
	history, err := pc.ProfileService.FetchSearchHistory(userpayload)
	if err != nil {
		if err.Error() == "no history found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "success",
				"message": "تاریخچه ای وجود ندارد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"history": history,
	})
}


func (pc *ProfileController) DeleteSearchHistory(c *gin.Context){
	userinfo, exists := c.Get("user")
	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)

	searchId := c.Param("id")
	err := pc.ProfileService.DeleteSearchHistory(userpayload, searchId)

	if err != nil{
		if err.Error() == "history not found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "success",
				"message": "تاریخچه ای وجود ندارد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
	}

	c.Status(http.StatusAccepted)
}

func (pc *ProfileController) ClearSearchHistory(c *gin.Context){
	userinfo, exists := c.Get("user")
	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)

	err := pc.ProfileService.ClearSearchHistory(userpayload)

	if err != nil{
		if err.Error() == "history not found"{
			c.JSON(http.StatusNotFound, gin.H{
				"status": "success",
				"message": "تاریخچه ای وجود ندارد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
	}

	c.Status(http.StatusAccepted)
}

func (pc *ProfileController) CreatePlaceLabel(c *gin.Context) {
	value, _ := c.Get("user")
	userinfo := value.(services_auth.UserInfo)

	var req serializers_profile.CreatePlaceLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ساختار داده ورودی اشتباه است.",
		})
		return
	}

	label, err := pc.ProfileService.CreatePlaceLabel(req.Name, userinfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}

	c.JSON(http.StatusCreated, label)
}

func (pc *ProfileController) GetUserPlaceLabels(c *gin.Context) {
	value, _ := c.Get("user")
	userinfo := value.(services_auth.UserInfo)
	labels, err := pc.ProfileService.GetUserPlaceLabels(userinfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "خطا در دریافت برچسب ها",
		})
		return
	}

	c.JSON(http.StatusOK, labels)
}

func (pc *ProfileController) DeletePlaceLabel(c *gin.Context) {
	value, _ := c.Get("user")
	userinfo := value.(services_auth.UserInfo)

	labelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "شناسه برچسب نامعتبر است",
		})
		return
	}

	err = pc.ProfileService.DeletePlaceLabel(uint(labelID), userinfo.ID)
	if err != nil {
		
		status := http.StatusInternalServerError
		message := "مشکل غیرمنتظره ای رخ داده است."

		if err.Error() == "record not found" {
			status = http.StatusNotFound
			message = "برچسب پیدا نشد"
		}

		c.JSON(status, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "",
	})
}

func (pc *ProfileController) CreateSavedLocation(c *gin.Context){
	var validateBody serializers_profile.CreateSavedLocationRequest
	if err := c.ShouldBindJSON(&validateBody); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "ساختار داده ورودی اشتباه است.",
		})
	}

	_, err := pc.PlaceService.GetPlaceByID(*validateBody.PlaceID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است.",
		})
	}

	value, _ := c.Get("user")
	userinfo := value.(services_auth.UserInfo)
	savedLocation, err := pc.ProfileService.CreateSavedLocation(userinfo.ID, services_profile.CreateSavedLocationParams{
		PlaceID: validateBody.PlaceID,
		PlaceLabelID: validateBody.PlaceLabelID,
		UserNote: validateBody.UserNote,
		Latitude: validateBody.Latitude,
		Longitude: validateBody.Longitude,
	})
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است.",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"saved_location": savedLocation,
		"status": "success",
		"message": "",
	})
}

func (pc *ProfileController) GetSavedLocationsByPlaceLabel(c *gin.Context){
	var validatedBody serializers_profile.GetSavedLocationsByPlaceLabelRequest
	if err := c.ShouldBindUri(&validatedBody); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "ساختار داده ورودی اشتباه است.",
		})
	}

	value, _ := c.Get("user")
	userinfo := value.(services_auth.UserInfo)

	places, err := pc.ProfileService.GetSavedLocationsByPlaceLabel(userinfo.ID, validatedBody.PlaceLabelID)
	if err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"message": "چنین برچسبی پیدا نشد.",
		})
		return 
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": places,
		"status": "success",
		"message": "",
	})
}

func (pc *ProfileController) UpdateSavedLocation(c *gin.Context) {
	var validatedBody serializers_profile.UpdateSavedLocationRequest
	if err := c.ShouldBindJSON(&validatedBody); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "ساختار داده ورودی اشتباه است",
		})
	}

	value, _ := c.Get("user")
	userinfo := value.(services_auth.UserInfo)

	savedLocation, err := pc.ProfileService.UpdateSavedLocation(
		services_profile.UpdateSavedLocationServiceParams{
			ID: validatedBody.ID,
			UserNote: validatedBody.UserNote,
			PlaceLabelID: validatedBody.PlaceLabelID,
			UserID: userinfo.ID,
		},
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"message": "مکان ذخیره شده یافت نشد.",
		})
		return 
	} else if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "خطا غیر منتظره ای رخ داده است.",
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"saved_location": savedLocation,
		"status": "success",
		"message": "",
	})
}

func (pc *ProfileController) HardDeleteSavedLocation(c *gin.Context){
	var validatedBody serializers_profile.DeleteSavedLocation
	if err := c.ShouldBindUri(&validatedBody); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "ساختار داده ورودی اشتباه هست.",
		})
		return 
	}

	value, _ := c.Get("user")
	userinfo := value.(services_auth.UserInfo)
	
	err := pc.ProfileService.HardDeleteSavedLocation(validatedBody.ID, userinfo.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"message": "مکان ذخیره شده یافت نشد.",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "خطای غیرمنتظره ای رخ داده است.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "",
	})
}
