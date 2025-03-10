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

	// Add each domain routes here
	authentication.AuthRoutes(router)

	router.Run(fmt.Sprintf("localhost:%s", initializers.SERVER_PORT))
}
