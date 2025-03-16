package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	"pashmak.com/pashmak/routers"
)

var (
	Router *gin.Engine
	DB *gorm.DB
	Redis *redis.Client
)

func init() {
	bootstrap.LoadEnvVars()
	DB = bootstrap.SetUpPostgres()
	Redis = bootstrap.SetupRedis()
	bootstrap.MakeMigrations(DB)
}

func main() {
	Router = gin.Default()
	
	// Add each domain routes here	
	routers_auth.AuthRoutes(Router, DB, Redis)

	Router.Run(fmt.Sprintf("localhost:%s", bootstrap.SERVER_PORT))
}