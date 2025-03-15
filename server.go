package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	"pashmak.com/pashmak/routers"
)

var (
	Router *gin.Engine
	DB *gorm.DB
)

func init() {
	bootstrap.LoadEnvVars()
	DB = bootstrap.SetUpPostgres()
	bootstrap.MakeMigrations(DB)
}

func main() {
	Router = gin.Default()
	
	// Add each domain routes here	
	routers_auth.AuthRoutes(Router, DB)

	Router.Run(fmt.Sprintf("localhost:%s", bootstrap.SERVER_PORT))
}