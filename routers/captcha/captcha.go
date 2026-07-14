package routers_captcha

import (
	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/bootstrap"
	controllers_captcha "pashmak.com/pashmak/controllers/captcha"
	services_captcha "pashmak.com/pashmak/services/captcha"
)

func CaptchaRoutes(router *gin.Engine, appConfig *bootstrap.AppConfig) {
	captchaService := services_captcha.NewCaptchaService(appConfig.ArcaptchaSiteKey, appConfig.ArcaptchaSecretKey)
	captchaController := controllers_captcha.NewCaptchaController(captchaService)
	captchaGroup := router.Group("/captcha")
	{
		captchaGroup.POST("/verify", captchaController.VerifyCaptcha)
	}
}
