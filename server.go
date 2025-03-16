package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	routers_auth "pashmak.com/pashmak/routers"
)

var (
	Router     *gin.Engine
	DB         *gorm.DB
	Redis      *redis.Client
)

func init() {
	bootstrap.LoadEnvVars()
	DB = bootstrap.SetUpPostgres()
	Redis = bootstrap.SetupRedis()
	bootstrap.MakeMigrations(DB)
}

func main() {
	Router = gin.Default()

	// FIXME: Should be checked if it's necessary in production or not
	// Global middleware to set CORS headers
	Router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		
	})

	// Add each domain routes here
	routers_auth.AuthRoutes(Router, DB, Redis)

	Router.Run(fmt.Sprintf("localhost:%s", bootstrap.SERVER_PORT))
}
