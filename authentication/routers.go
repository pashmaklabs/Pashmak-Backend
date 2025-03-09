package authentication

import (
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	//Add routes here
	router.POST("StartEmailAuth", StartEmailAuth)
}
