package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/initializers"
)

func init() {
	initializers.LoadEnvVars()
	db := initializers.SetUpPostgres()
	initializers.MakeMigrations(db)
}

func main() {
	router := gin.Default()
	router.Run(fmt.Sprintf("localhost:%s", initializers.SERVER_PORT))
}