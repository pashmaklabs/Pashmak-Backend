package middlewares_cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/bootstrap"
)

type CorsMiddleware struct {
	appConfig *bootstrap.AppConfig
}

func NewCorsMiddleware(appConfig *bootstrap.AppConfig) *CorsMiddleware {
	return &CorsMiddleware{appConfig: appConfig}
}

func (cm *CorsMiddleware) SetCORSHeader() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowOrigins = cm.appConfig.AllowdOrigins
    config.AllowCredentials = true
    config.AllowHeaders = append(config.AllowHeaders, "Authorization") // Add any additional headers your frontend sends
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    return cors.New(config)
}