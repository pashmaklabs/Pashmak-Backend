package authentication

import (
	"github.com/gin-gonic/gin"
)

var(
	routeService = NewAuthService()
	routeController = NewAuthController(routeService)
)

func AuthRoutes(router *gin.Engine) {
	
	
	//Add routes here
	router.POST("StartEmailAuth", routeController.StartEmailAuth)
}
