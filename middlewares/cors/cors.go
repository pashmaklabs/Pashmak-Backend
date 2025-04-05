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
	return cors.New(config)
}