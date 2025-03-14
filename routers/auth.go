package routers_auth


import (
	"pashmak.com/pashmak/services"
	"pashmak.com/pashmak/controllers"
	"github.com/gin-gonic/gin"

)

var(
	routeService = services_auth.NewAuthService()
	routeController = controlllers_auth.NewAuthController(routeService)
)

func AuthRoutes(router *gin.Engine) {
	//Add routes here
	router.POST("StartEmailAuth", routeController.StartEmailAuth)
}
