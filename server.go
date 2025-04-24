package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	middlewares_cors "pashmak.com/pashmak/middlewares/cors"
	routers_auth "pashmak.com/pashmak/routers/auth"
	routers_comment "pashmak.com/pashmak/routers/comment"
	routers_navigation "pashmak.com/pashmak/routers/navigation"
	routers_profile "pashmak.com/pashmak/routers/profile"
	routers_place "pashmak.com/pashmak/routers/place"
	"pashmak.com/pashmak/serializers"
)

var (
	Router    *gin.Engine
	DB        *gorm.DB
	Redis     *redis.Client
	Minio     *minio.Client
	AppConfig *bootstrap.AppConfig
)

func init() {
	AppConfig = bootstrap.LoadEnvVars()
	DB = bootstrap.SetUpPostgres()
	Redis = bootstrap.SetupRedis()
	Minio = bootstrap.SetUpMinio()
	bootstrap.MakeMigrations(DB)
}

func main() {
	Router = gin.Default()

	corsMiddleware := middlewares_cors.NewCorsMiddleware(AppConfig)
	Router.Use(corsMiddleware.SetCORSHeader())

	// Set up the Gin validator with custom validation rules
	validate := validator.New()
	serializers.RegisterCustomValidators(validate)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		serializers.RegisterCustomValidators(v)
	} else {
		log.Println("Validator engine cannot be cast to validator.Validate")
	}

	// Add each domain routes here
	routers_auth.AuthRoutes(Router, DB, Redis, AppConfig)
	routers_profile.ProfileRoutes(Router, DB, Redis, AppConfig)
	routers_navigation.NavigationRoutes(Router, DB, AppConfig)
	routers_comment.CommentRoutes(Router, DB, Redis, AppConfig)
	routers_place.PlaceRoutes(Router, DB, Redis, AppConfig)
	
	Router.Run(fmt.Sprintf(":%s", AppConfig.ServerPort))
}
