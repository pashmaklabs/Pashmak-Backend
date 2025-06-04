package controllers_place

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	serializers_place "pashmak.com/pashmak/serializers/place"
	services_auth "pashmak.com/pashmak/services/auth"
	services_comment "pashmak.com/pashmak/services/comment"
	services_place "pashmak.com/pashmak/services/place"
)

type PlaceController struct {
	PlaceService   *services_place.PlaceService
	CommnetService *services_comment.CommentService
}

func NewPlaceController(placeService *services_place.PlaceService, commentService *services_comment.CommentService) *PlaceController {
	return &PlaceController{
		PlaceService:   placeService,
		CommnetService: commentService,
	}
}

func (pc *PlaceController) GetPlace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "شناسه مکان نامعتبر است",
		})
		return
	}

	place, err := pc.PlaceService.GetPlaceByID(uint(id))
	if err != nil {
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
		return
	}

	avgRating, err := pc.CommnetService.GetAverageRating(strconv.FormatUint(id, 10))
	if err != nil {
		if err.Error() == "نظری ثبت نشده" {
			avgRating = 0
		} else {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "خطا در دریافت امتیاز",
			})
			return
		}
	}

	response := serializers_place.PlaceWithRatingResponse{
		Name:      place.Name,
		Amenity:   place.Amenity,
		Latitude:  place.Latitude,
		Longitude: place.Longitude,
		Rating:    avgRating,
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "",
		"place":   response,
	})
}

func (pc *PlaceController) SearchPlace(c *gin.Context) {
	q := c.Query("q")
	lat := c.Query("lat")
	long := c.Query("lng")
	query := fmt.Sprintf("Query: %s\nLatitude: %s\nLongitude: %s", q, lat, long)
	places, err := pc.PlaceService.SearchPlace(query)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "خطا در جستجوی مکان",
		})
		return
	}


	sessionID, err := c.Cookie("session_id")
    if err != nil || sessionID == "" {
        sessionID = uuid.New().String()
        c.SetCookie("session_id", sessionID, 30*24*3600, "/", "", false, true) // consider
    }
	userinfo, loggedIn := c.Get("user")
	if loggedIn{
		userpayload := userinfo.(services_auth.UserInfo)
		log.Println("v ...any")
		// Get user_id if logged in
		var userID uint
		userID = userinfo.(services_auth.UserInfo).ID // Adjust to your UserInfo struct
		
		err = pc.PlaceService.SaveSearch(userID, sessionID, loggedIn, q)
		if err != nil{
			log.Printf("Failed to save search query fo user: %v, %v", userpayload.ID, err)
		}
	}
	
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "",
		"places":  places,
	})
}

