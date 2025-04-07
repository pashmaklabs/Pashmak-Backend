package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	middlewares_cors "pashmak.com/pashmak/middlewares/cors"
	routers_auth "pashmak.com/pashmak/routers"
)

var (
	Router    *gin.Engine
	DB        *gorm.DB
	Redis     *redis.Client
	AppConfig *bootstrap.AppConfig
)

func init() {
	AppConfig = bootstrap.LoadEnvVars()
	DB = bootstrap.SetUpPostgres()
	Redis = bootstrap.SetupRedis()
	bootstrap.MakeMigrations(DB)
}

func main() {
	Router = gin.Default()
	log.Println("+", "password", "+")
	corsMiddleware := middlewares_cors.NewCorsMiddleware(AppConfig)
	Router.Use(corsMiddleware.SetCORSHeader())

	// Add each domain routes here
	routers_auth.AuthRoutes(Router, DB, Redis, AppConfig)

	Router.Run(fmt.Sprintf(":%s", AppConfig.ServerPort))
}
