package controllers_captcha

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pashmak.com/pashmak/serializers/captcha"
	services_captcha "pashmak.com/pashmak/services/captcha"
)

type CaptchaController struct {
	service *services_captcha.CaptchaService
}

func NewCaptchaController(captchaservice *services_captcha.CaptchaService) *CaptchaController {
	return &CaptchaController{
		service: captchaservice,
	}
}

func (c *CaptchaController) VerifyCaptcha(ctx *gin.Context) {
	var req captcha.VerificationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := c.service.VerifyToken(req.Token); err != nil {
		ctx.JSON(http.StatusBadRequest, captcha.VerificationResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// You can add additional logic here, like storing the verified status in context
	ctx.Set("captcha_verified", true)

	ctx.JSON(http.StatusOK, captcha.VerificationResponse{
		Success: true,
		Message: "Captcha verified successfully",
	})
}

// Middleware to verify captcha on specific routes
// func (c *CaptchaController) CaptchaMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		token := ctx.GetHeader("X-Captcha-Token")
// 		if token == "" {
// 			token = ctx.Query("captcha_token")
// 		}

// 		if token == "" {
// 			ctx.JSON(http.StatusBadRequest, gin.H{
// 				"success": false,
// 				"message": "Captcha token required",
// 			})
// 			ctx.Abort()
// 			return
// 		}

// 		if err := c.service.VerifyToken(token); err != nil {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{
// 				"success": false,
// 				"message": "Invalid captcha",
// 			})
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set("captcha_verified", true)
// 		ctx.Next()
// 	}
// }
