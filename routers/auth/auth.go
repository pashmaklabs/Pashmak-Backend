package routers_auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	controllers_auth "pashmak.com/pashmak/controllers/auth"
	middlewares_auth "pashmak.com/pashmak/middlewares/auth"
	middlewares_ratelimit "pashmak.com/pashmak/middlewares/ratelimit"
	middlewares_validation "pashmak.com/pashmak/middlewares/validation"
	serializers_auth "pashmak.com/pashmak/serializers/auth"
	services_auth "pashmak.com/pashmak/services/auth"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client, appConfig *bootstrap.AppConfig) {
	authService := services_auth.NewAuthService(db, redisClient, appConfig)
	authController := controllers_auth.NewAuthController(authService, appConfig)
	authMiddleware := middlewares_auth.NewAuthMiddleware(authService)

	// OTP send: 5 requests / 10 minutes per IP — most sensitive, sends email
	otpSendLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 5, 10*time.Minute, "otp_send", middlewares_ratelimit.KeyByIP)

	// OTP verify: 10 attempts / 5 minutes per IP — brute-force on 4-digit OTP
	otpVerifyLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 10, 5*time.Minute, "otp_verify", middlewares_ratelimit.KeyByIP)

	// Password login: 10 attempts / 5 minutes per IP
	loginLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 10, 5*time.Minute, "login", middlewares_ratelimit.KeyByIP)

	// Forget password send: 3 / 15 minutes — prevents email spam
	forgetSendLimiter := middlewares_ratelimit.NewRateLimiter(redisClient, 3, 15*time.Minute, "forget_send", middlewares_ratelimit.KeyByIP)

	auth := router.Group("/auth")
	{
		auth.POST("/otp/send",
			otpSendLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_auth.SendOTPRequest](),
			authController.SendOTP)

		auth.POST("/otp/verify",
			otpVerifyLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_auth.VerifyOTPRequest](),
			authController.VerifyOTP)

		auth.GET("/protected", authMiddleware.LoginMiddleware(), authController.ProtectedRouter)

		auth.POST("/password",
			loginLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_auth.LoginWithPasswordRequest](),
			authController.LoginWithPassword)

		auth.POST("/password/forget/send",
			forgetSendLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_auth.ForgetPasswordRequest](),
			authController.ForgetPassword)

		auth.POST("/password/forget/verify",
			otpVerifyLimiter.Middleware(),
			middlewares_validation.ValidationMiddleware[serializers_auth.ForgetPasswordVerifyRequest](),
			authController.ForgetPasswordVerify)

		auth.POST("/password/forget/reset",
			authMiddleware.LoginMiddleware(),
			middlewares_validation.ValidationMiddleware[serializers_auth.ForgetPasswordResetRequest](),
			authController.ForgetPasswordReset)

		auth.PATCH("/signup",
			authMiddleware.LoginMiddleware(),
			middlewares_validation.ValidationMiddleware[serializers_auth.SignUpRequest](),
			authController.SignUp)
	}
}