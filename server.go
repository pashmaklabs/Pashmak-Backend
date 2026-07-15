package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
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
	routers_captcha "pashmak.com/pashmak/routers/captcha"
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
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              AppConfig.SentryDsn,
		AttachStacktrace: true,
		Environment:      AppConfig.Environment,
		TracesSampler: sentry.TracesSampler(func(ctx sentry.SamplingContext) float64 {
			// Don't trace health checks
			if strings.Contains(ctx.Span.Name, "/health") {
				return 0
			}

			// Always trace critical endpoints (payments, auth)
			if strings.Contains(ctx.Span.Name, "/payment") ||
				strings.Contains(ctx.Span.Name, "/login") {
				return 1.0
			}

			// Sample everything else at 10% for production
			return 0.1
		}),
	})

	if err != nil {
		log.Fatalf("sentry.Init failed: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

	Router = gin.Default()

	// Initialize Prometheus metrics first
	middlewares_prometheus.PrometheusInit()

	// Apply Prometheus middleware to all routes (should be early in the chain)
	Router.Use(middlewares_prometheus.PrometheusMiddleware())

	// Apply Sentry Middleware right after Prometheus
	Router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

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

	// --- SENTRY TEST ROUTES (Remove these before deploying to production) ---
	Router.GET("/sentry-test-error", func(c *gin.Context) {
		// 1. Testing a manually captured, handled error
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.CaptureMessage("Sentry manual test: Something minor went wrong!")
		}
		c.JSON(200, gin.H{"status": "Handled message sent to Sentry"})
	})

	Router.GET("/sentry-test-panic", func(c *gin.Context) {
		// 2. Testing an unhandled panic (Nil pointer dereference)
		// Sentry middleware will automatically catch this, report it, and keep the server alive!
		var emptyStringPointer *string
		println(*emptyStringPointer) // This intentionally crashes the request
	})

	// Expose Prometheus metrics endpoint
	Router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Add each domain routes here
	routers_auth.AuthRoutes(Router, DB, Redis, AppConfig)
	routers_profile.ProfileRoutes(Router, DB, PGVectorDB, Redis, Minio, AppConfig)
	routers_navigation.NavigationRoutes(Router, DB, Redis, AppConfig)
	routers_comment.CommentRoutes(Router, DB, PGVectorDB, Redis, AppConfig)
	routers_place.PlaceRoutes(Router, DB, PGVectorDB, Redis, Minio, AppConfig)
	routers_admin.AdminRoutes(Router, DB, PGVectorDB, Redis, Minio, AppConfig)
	routers_captcha.CaptchaRoutes(Router, AppConfig)

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
