package routers_navigation

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_navigation "pashmak.com/pashmak/controllers/navigation"
	services_navigation "pashmak.com/pashmak/services/navigation"
)

func NavigationRoutes(c *gin.Engine, db *gorm.DB, appconfig *bootstrap.AppConfig) {
	navigationService := services_navigation.NewNavigationService(db, appconfig)
	navigationController := controllers_navigation.NewNavigationController(navigationService, db, appconfig)
		
	navigation := c.Group("/navigation")
	{
		navigation.GET("/route", navigationController.GetRoute)
	}
}
