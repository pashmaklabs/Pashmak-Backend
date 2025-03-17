package routers_auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	controlllers_auth "pashmak.com/pashmak/controllers"
	services_auth "pashmak.com/pashmak/services/auth"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB) {
	routeService := services_auth.NewAuthService(db)
	routeController := controlllers_auth.NewAuthController(routeService)
	
	auth := router.Group("/auth")
	{
		auth.POST("/send-otp", routeController.SendOTP)
	}
}