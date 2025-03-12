package authentication

import (
	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/initializers"
)

var (
	authService    = NewAuthService(initializers.DB, initializers.RedisClient)
	authController = NewAuthController(authService)
)

func AuthRoutes(router *gin.Engine) {
	//Add routes here
	router.POST("/StartEmailAuth", authController.StartEmailAuth)
	router.POST("/VerifyOTP", authController.VerifyOTP)
}
