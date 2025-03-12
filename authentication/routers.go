package authentication

import (
	"log"
	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/initializers"
)

var (
	authService    = NewAuthService(initializers.DB, initializers.RedisClient)
	authController = NewAuthController(authService)
)

// func init() {
// 	log.Println("Initializing AuthService...")

// 	if authService.DB == nil {
// 		log.Println("Warning: authService is nil")
// 		log.Println(initializers.DB)
// 		log.Println(initializers.RedisClient)
// 		panic("err")
// 	}

// }

func AuthRoutes(router *gin.Engine) {
	//Add routes here
	log.Println(initializers.DB)
	log.Println(initializers.RedisClient)
	router.POST("/StartEmailAuth", authController.StartEmailAuth)
}
