package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/bootstrap"
	"pashmak.com/pashmak/routers"
)

var (
	router *gin.Engine
)

func init() {
	bootstrap.LoadEnvVars()
	db := bootstrap.SetUpPostgres()
	redis_client := bootstrap.SetupRedis()
	bootstrap.MakeMigrations(db)
	router = gin.Default()
	
	// Add each domain routes here	
	routers_auth.AuthRoutes(router, db, redis_client)
}

func main() {
	router.Run(fmt.Sprintf("localhost:%s", bootstrap.SERVER_PORT))
}