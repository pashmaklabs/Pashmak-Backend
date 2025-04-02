package middlewares_cors

import (
	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/bootstrap"
)

type CorsMiddleware struct {
	appConfig *bootstrap.AppConfig
}

func NewCorsMiddleware(appConfig *bootstrap.AppConfig) *CorsMiddleware {
	return &CorsMiddleware{appConfig: appConfig}
}

func (cm *CorsMiddleware)SetCORSHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cm.appConfig.AllowdOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}