package routers_auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pashmak.com/pashmak/controllers"
	"pashmak.com/pashmak/services"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB) {
	//Add routes here
	routeService := services_auth.NewAuthService(db)
	routeController := controlllers_auth.NewAuthController(routeService)
	router.POST("StartEmailAuth", routeController.StartEmailAuth)
}
