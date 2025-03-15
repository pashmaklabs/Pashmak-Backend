package routers_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/controllers"
	"pashmak.com/pashmak/services"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client) {
	routeService := services_auth.NewAuthService(db, redis)
	routeController := controlllers_auth.NewAuthController(routeService)
	
	//Add routes here
	router.POST("/StartEmailAuth", routeController.StartEmailAuth)
}