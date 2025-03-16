package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/authentication"
	"pashmak.com/pashmak/initializers"
	
)

func init() {
	initializers.LoadEnvVars()
	db := initializers.SetUpPostgres()
	initializers.MakeMigrations(db)
	initializers.SetupRedis()
}

func main() {
	router := gin.Default()

	// Global middleware to set CORS headers
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add each domain routes here
	authentication.AddAuthRoutes(router)

	router.Run(fmt.Sprintf("localhost:%s", initializers.SERVER_PORT))
}
