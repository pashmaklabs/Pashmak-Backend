package controlllers_auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"errors"

	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	serializers_auth "pashmak.com/pashmak/serializers"
	services_auth "pashmak.com/pashmak/services/auth"
)

type AuthController struct {
	authService *services_auth.AuthService
	AppConfig 	*bootstrap.AppConfig
}

func NewAuthController(authService *services_auth.AuthService, appConfig *bootstrap.AppConfig) *AuthController {
	return &AuthController{
		authService: authService,
		AppConfig: appConfig,
	}
}

func (ac *AuthController) SendOTP(c *gin.Context) {
	// Read body
	var body serializers_auth.SendOTPRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    "error",
			"message":   "در خواندن بدنه ی درخواست خطایی رخ داد",
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
			"message":   "رمز یکبار مصرف ارسال شد",
			"userExists":    false,
		})
		return
	} else{
		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"message":   "رمز یکبار مصرف ارسال شد",
			"userExists":    true,
		})
		return
	}
}

func (ac *AuthController) VerifyOTP(c *gin.Context) {
	var body serializers_auth.VerifyOTPRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    "error",
			"message":   "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}

	resp, err := ac.authService.ValidateOTP(body.Email, body.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		return
	}
	if resp {
		user, err := ac.authService.GetUserByGmail(body.Email)
		exists := true
		if errors.Is(err, gorm.ErrRecordNotFound) {
			exists = false
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "مشکل غیرمنتظره ای رخ داده است",
			})
			return
		}
		// TODO: Move logic to service
		jwt, err := ac.authService.GenerateJWT(user)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "مشکل غیرمنتظره ای رخ داده است",
			})
			return
		}
		c.SetCookie("jwt_token", jwt, int(ac.AppConfig.TokenAge), "/", "", false, true)
		if !exists{
			err := ac.authService.CreateUser(body.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "مشکل غیرمنتظره ای رخ داده است",
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "ورود با موفقیت انجام شد.",
		})
		return
	}else{
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "رمز یکبار مصرف اشتباه وارد شده.",
		})
		return
	}
}

func (ac *AuthController) ProtectedRouter(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "این یک api محافظت شده است :)",
	})
	return
}

func (ac *AuthController) SignUp(c *gin.Context) {
	var body serializers_auth.SignUpRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    "error",
			"message":   "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}
	userinfo, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	err := ac.authService.SignUp(userinfo.(services_auth.UserInfo).Email, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیر منتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "ثبت نام با موفقیت انجام شد.",
	})
	return
}