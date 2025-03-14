package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/bootstrap"
	"pashmak.com/pashmak/routers"
)

func init() {
	bootstrap.LoadEnvVars()
	db := bootstrap.SetUpPostgres()
	bootstrap.MakeMigrations(db)
}

func main() {
	router := gin.Default()

	// Add each domain routes here
	routers_auth.AuthRoutes(router)

	router.Run(fmt.Sprintf("localhost:%s", bootstrap.SERVER_PORT))
}
