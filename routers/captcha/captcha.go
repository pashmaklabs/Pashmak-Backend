package routers_captcha

import (
    "github.com/gin-gonic/gin"
    "pashmak.com/pashmak/controllers/captcha"
)

func CaptchaRoutes(router *gin.RouterGroup, controller *captcha.CaptchaController) {
    captchaGroup := router.Group("/captcha")
    {
        captchaGroup.POST("/verify", controller.VerifyCaptcha)
    }
}