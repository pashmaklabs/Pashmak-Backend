package routers_auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_auth "pashmak.com/pashmak/controllers/auth"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	middlewares_validation "pashmak.com/pashmak/middlewares/validation"
	serializers_auth "pashmak.com/pashmak/serializers/auth"
	services_auth "pashmak.com/pashmak/services/auth"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client, appConfig *bootstrap.AppConfig) {
	authService := services_auth.NewAuthService(db, redis, appConfig)
	authController := controllers_auth.NewAuthController(authService, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)
	
	auth := router.Group("/auth")
	{
		auth.POST("/otp/send",
			middlewares_validation.ValidationMiddleware(serializers_auth.SendOTPRequest{}),
			authController.SendOTP)
		auth.POST("/otp/verify",
			middlewares_validation.ValidationMiddleware(serializers_auth.VerifyOTPRequest{}),
			authController.VerifyOTP)
		auth.GET("/protected", authMiddleware.LoginMiddleware(), authController.ProtectedRouter)
		auth.POST("/password", 
			middlewares_validation.ValidationMiddleware(serializers_auth.LoginWithPasswordRequest{}),
			authController.LoginWithPassword)
		auth.POST("/password/forget/send",
			middlewares_validation.ValidationMiddleware(serializers_auth.SendOTPRequest{}),
			authController.ForgetPassword)
		auth.POST("/password/forget/verify",
			middlewares_validation.ValidationMiddleware(serializers_auth.VerifyOTPRequest{}),
			authController.ForgetPasswordVerify)
		auth.POST("/password/forget/reset",
			authMiddleware.LoginMiddleware(), 
			middlewares_validation.ValidationMiddleware(serializers_auth.ForgetPasswordResetRequest{}),
			authController.ForgetPasswordReset)
		// TODO: Why not put?
		auth.PATCH("/signup",
			authMiddleware.LoginMiddleware(),
			middlewares_validation.ValidationMiddleware(serializers_auth.SignUpRequest{}),
            authController.SignUp,
        )
	}
}