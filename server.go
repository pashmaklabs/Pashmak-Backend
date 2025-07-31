package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"pashmak.com/pashmak/bootstrap"
	middlewares_cors "pashmak.com/pashmak/middlewares/cors"
	middlewares_prometheus "pashmak.com/pashmak/middlewares/prometheus"
	routers_admin "pashmak.com/pashmak/routers/admin"
	routers_auth "pashmak.com/pashmak/routers/auth"
	routers_comment "pashmak.com/pashmak/routers/comment"
	routers_navigation "pashmak.com/pashmak/routers/navigation"
	routers_place "pashmak.com/pashmak/routers/place"
	routers_profile "pashmak.com/pashmak/routers/profile"
	"pashmak.com/pashmak/serializers"
)

var (
	Router     *gin.Engine
	DB         *gorm.DB
	PGVectorDB *gorm.DB
	Redis      *redis.Client
	Minio      *minio.Client
	AppConfig  *bootstrap.AppConfig
)

func init() {
	AppConfig = bootstrap.LoadEnvVars()
	DB = bootstrap.SetUpPostgres()
	PGVectorDB = bootstrap.SetUpPGVector()
	Redis = bootstrap.SetupRedis()
	Minio = bootstrap.SetUpMinio()
	bootstrap.MakeMigrations(DB)
	bootstrap.MakePGVectorMigrations(PGVectorDB)
}

func main() {
	Router = gin.Default()

	// Initialize Prometheus metrics first
	middlewares_prometheus.PrometheusInit()

	// Apply Prometheus middleware to all routes (should be early in the chain)
	Router.Use(middlewares_prometheus.PrometheusMiddleware())

	// CORS middleware
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

	// Health check endpoint
	Router.GET("/health", middlewares_prometheus.HealthHandler)

	// Expose Prometheus metrics endpoint
	Router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Add each domain routes here
	routers_auth.AuthRoutes(Router, DB, Redis, AppConfig)
	routers_profile.ProfileRoutes(Router, DB, PGVectorDB, Redis, Minio, AppConfig)
	routers_navigation.NavigationRoutes(Router, DB, AppConfig)
	routers_comment.CommentRoutes(Router, DB, PGVectorDB, Redis, AppConfig)
	routers_place.PlaceRoutes(Router, DB, PGVectorDB, Redis, Minio, AppConfig)
	routers_admin.AdminRoutes(Router, DB, PGVectorDB, Redis, Minio, AppConfig)

	// Start periodic database metrics collection
	go updateDatabaseMetrics()

	// Start server
	Router.Run(fmt.Sprintf(":%s", AppConfig.ServerPort))
}

// Periodically update database connection metrics
func updateDatabaseMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if sqlDB, err := DB.DB(); err == nil {
				stats := sqlDB.Stats()
				middlewares_prometheus.UpdateDatabaseMetrics(
					stats.OpenConnections,
					stats.Idle,
				)
			}

			// Update Redis metrics if possible
			if Redis != nil {
				// Redis client doesn't expose connection pool stats directly
				// You might need to implement this based on your Redis setup
				middlewares_prometheus.UpdateRedisMetrics(1) // placeholder
			}
		}
	}
}