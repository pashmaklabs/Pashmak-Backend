package controlllers_auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serializers_auth "pashmak.com/pashmak/serializers"
	services_auth "pashmak.com/pashmak/services"
)

type AuthController struct {
	authService *services_auth.AuthService
}

func NewAuthController(authService *services_auth.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) SendOTP(c *gin.Context) {
	// Read body
	var body serializers_auth.SendOTPRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    "error",
			"message":   "در خواندن بدنه ی درخواست خطایی رخ داد",
			"errorCode": "INVALID_REQUEST_BODY",
		})
		return
	}

	// Pass to auth service
	resp, err := ac.authService.ValidateUser(body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if !resp {
		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"message":   "کاربر یافت نشد",
			"exists":    false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "کاربر یافت شد",
		"exists":    true,
	})

}

func (ac *AuthController) VerifyOTP(c *gin.Context) {
	var body serializers_auth.VerifyOTPRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    "error",
			"message":   "در خواندن بدنه ی درخواست خطایی رخ داد",
			"errorCode": "INVALID_REQUEST_BODY",
		})
		return
	}

	resp, err := ac.authService.ValidateOTP(body.Email, body.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if !resp {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "OTP mismatch",
		})
	}
	if resp {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "OTP match",
		})
	}
}
