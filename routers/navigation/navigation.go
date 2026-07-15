package routers_navigation

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_navigation "pashmak.com/pashmak/controllers/navigation"
	middlewares_ratelimit "pashmak.com/pashmak/middlewares/ratelimit"
	services_navigation "pashmak.com/pashmak/services/navigation"
)

func NavigationRoutes(c *gin.Engine, db *gorm.DB, redisClient *redis.Client, appconfig *bootstrap.AppConfig) {
	navigationService := services_navigation.NewNavigationService(db, appconfig)
	navigationController := controllers_navigation.NewNavigationController(navigationService, db, appconfig)

	// Route calculation is likely expensive (external API call or graph traversal)
	routeLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 30, time.Minute, "navigation_route", middlewares_ratelimit.KeyByIP)

	navigation := c.Group("/navigation")
	{
		navigation.GET("/route", routeLimiter.Middleware(), navigationController.GetRoute)
	}
}