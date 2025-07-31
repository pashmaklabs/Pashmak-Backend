package controllers_place

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	models_place "pashmak.com/pashmak/models/place"
	serializers_place "pashmak.com/pashmak/serializers/place"
	services_auth "pashmak.com/pashmak/services/auth"
	services_comment "pashmak.com/pashmak/services/comment"
	services_place "pashmak.com/pashmak/services/place"
	"pashmak.com/pashmak/services/placeOsmUtils"
	services_profile "pashmak.com/pashmak/services/profile"
)

type PlaceController struct {
	PlaceService   *services_place.PlaceService
	CommnetService *services_comment.CommentService
	ProfileService *services_profile.ProfileService
	AppConfig      *bootstrap.AppConfig
}

func NewPlaceController(placeService *services_place.PlaceService, commentService *services_comment.CommentService, profileService *services_profile.ProfileService, appConfig *bootstrap.AppConfig) *PlaceController {
	return &PlaceController{
		PlaceService:   placeService,
		CommnetService: commentService,
		ProfileService: profileService,
		AppConfig:      appConfig,
	}
}

func (pc *PlaceController) GetPlace(c *gin.Context) {
	idStr := c.Param("id")

	place, err := pc.PlaceService.GetPlaceByID(idStr)
	if err != nil {
		if strings.Contains(err.Error(), "invalid place ID format") {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "شناسه مکان نامعتبر است",
			})
			return
		}
		if strings.Contains(err.Error(), "no place found") {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "مکان یافت نشد",
			})
			return
		}
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "خطا در دریافت مکان",
		})
		log.Println(err)
		return
	}
	_, err = strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		place.Rating, err = pc.CommnetService.GetAverageRating(idStr)
		if err != nil {
			if err.Error() == "نظری ثبت نشده" {
				place.Rating = 0
			} else {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"status":  "error",
					"message": "خطا در دریافت امتیاز",
				})
				return
			}
		}
	}
	

	response := serializers_place.PlaceWithRatingResponse{
		ID:        place.ID, // Now this is a string
		Name:      place.Name,
		Amenity:   *place.Amenity,
		Latitude:  *place.Latitude,
		Longitude: *place.Longitude,
		Rating:    place.Rating,
		ImageURLs: place.ImageURLs,
	}

	value, exists := c.Get("user")
	if exists {
		userinfo := value.(services_auth.UserInfo)
		// Convert string ID to uint for the profile service
		placeIDUint, err := strconv.ParseUint(place.ID, 10, 32)
		if err == nil {
			savedLocation, err := pc.ProfileService.GetLabelOfPlace(userinfo.ID, uint(placeIDUint))
			if err == nil && savedLocation != nil {
				response.SavedLocation = savedLocation
			}
		}
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"place":   response,
	})
}

func (pc *PlaceController) SearchPlace(c *gin.Context) {
	q := c.Query("q")
	lat := c.Query("lat")
	long := c.Query("lng")
	agentic := c.Query("agentic") == "true"

	sessionID, err := c.Cookie("session_id")
	userinfo, loggedIn := c.Get("user")
	if !loggedIn && (err != nil || sessionID == "") {
		sessionID = uuid.New().String()
		c.SetCookie("session_id", sessionID, 30*24*3600, "/", pc.AppConfig.CookieDomain, false, true)
	}

	var userID *uint
	if loggedIn {
		id := userinfo.(services_auth.UserInfo).ID
		userID = &id
	}
	err = pc.PlaceService.SaveSearch(userID, sessionID, loggedIn, q)
	if err != nil {
		log.Printf("Failed to save search query fo user: %v, %v", userID, err)
	}

	places, err := pc.PlaceService.SearchPlace(q, lat, long, agentic)
	if err != nil {
		fmt.Println("err", err)
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "خطا در جستجوی مکان",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "",
		"places":  places,
	})
}

func (pc *PlaceController) GetPlaceRecommendations(c *gin.Context) {
	query := c.Query("q")

	sessionID, err := c.Cookie("session_id")
	userinfo, loggedIn := c.Get("user")
	if !loggedIn && (err != nil || sessionID == "") {
		sessionID = uuid.New().String()
		c.SetCookie("session_id", sessionID, 30*24*3600, "/", pc.AppConfig.CookieDomain, false, true)
	}

	var userID *uint
	if loggedIn {
		id := userinfo.(services_auth.UserInfo).ID
		userID = &id
	}
	err = pc.PlaceService.SaveSearch(userID, sessionID, loggedIn, query)
	if err != nil {
		log.Printf("Failed to save search query fo user: %v, %v", userID, err)
	}

	recommendations, err := pc.PlaceService.GetPlaceRecommendations(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "مشکل غیرمنتظره ای پیش آمده است"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"message":         "",
		"recommendations": recommendations,
	})
}

// UploadPlaceImage handles POST /places/:id/images for uploading a new image to a place.
func (pc *PlaceController) UploadPlaceImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "شناسه مکان نامعتبر است"})
		return
	}

	var place models_place.Place
	err = pc.PlaceService.DB.First(&place, uint(id)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			place, err := placeOsmUtils.ImportFromOSM(uint(id), pc.CommnetService.DB)
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
			}
			// Retrieve the newly created place
			if err := pc.PlaceService.DB.Where("id = ?", id).First(&place).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "خطا در دریافت رکورد مکان"})
			return
		}
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "فایل ارسال نشد"})
		log.Println(err)
		return
	}
	objectName, err := pc.PlaceService.UploadPlaceImage(&place, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "مشکل غیر منتظره ای رخ داده است"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "File uploaded successfully",
		"data": map[string]interface{}{
			"objectName": objectName,
		},
	})
}

// GetPlaceImage handles GET /places/:id/images/:image_name for retrieving a place image.
func (pc *PlaceController) GetPlaceImage(c *gin.Context) {
	idStr := c.Param("id")
	imageName := c.Param("image_name")
	_, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "شناسه مکان نامعتبر است"})
		return
	}
	imgStream, eTag, err := pc.PlaceService.GetPlaceImage(0, imageName) // placeID not used in service
	if err != nil {
		if err.Error() == "image not found" {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	defer imgStream.Close()
	c.Header("ETag", eTag)
	c.Header("Cache-Control", "public, max-age=3600")
	c.Header("Content-Type", "image/webp")
	_, err = io.Copy(c.Writer, imgStream)
	if err != nil {
		log.Println("Error writing image using Copy func:", err)
	}
}

func (pc *PlaceController) AddNewPlace(c *gin.Context) {
	validatedData, exists := c.Get("validated")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed to retrieve validated data from context: exists=%v", exists)
		return
	}

	body, ok := validatedData.(serializers_place.AddPlaceRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Printf("Failed type assertion for validated data: expected AddCommentRequest, got %T", validatedData)
		return
	}

	userinfo, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "ابتدا باید وارد شوید",
		})
		return
	}
	userpayload := userinfo.(services_auth.UserInfo)
	err := pc.PlaceService.AddNewPlace(userpayload, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیر منتظره ای رخ داد",
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "مکان با موفقیت ثبت شد",
	})
}
