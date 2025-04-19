package controllers_navigation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	services_navigation "pashmak.com/pashmak/services/navigation"
)

type NavigationController struct {
	NavigationService *services_navigation.NavigationService
	DB                *gorm.DB
	AppConfig         *bootstrap.AppConfig
}

func NewNavigationController(navigationService *services_navigation.NavigationService, db *gorm.DB, appConfig *bootstrap.AppConfig) *NavigationController {
	return &NavigationController{
		NavigationService: navigationService,
		DB:                db,
		AppConfig:         appConfig,
	}
}

func (rc *NavigationController) GetRoute(c *gin.Context) {
	startLat := c.Query("start_lat")
	startLon := c.Query("start_lon")
	endLat := c.Query("end_lat")
	endLon := c.Query("end_lon")

	if startLat == "" || startLon == "" || endLat == "" || endLon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	route, err := rc.NavigationService.FetchRoute(startLat, startLon, endLat, endLon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}
