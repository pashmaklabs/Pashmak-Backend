package controllers_auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"errors"

	"gorm.io/gorm"
	"pashmak.com/pashmak/bootstrap"
	serializers_auth "pashmak.com/pashmak/serializers/auth"
	services_auth "pashmak.com/pashmak/services/auth"
)

type AuthController struct {
	authService *services_auth.AuthService
	AppConfig   *bootstrap.AppConfig
}

func NewAuthController(authService *services_auth.AuthService, appConfig *bootstrap.AppConfig) *AuthController {
	return &AuthController{
		authService: authService,
		AppConfig:   appConfig,
	}
}

func (ac *AuthController) SendOTP(c *gin.Context) {
	var body serializers_auth.SendOTPRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}

	resp, err := ac.authService.ValidateUser(body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "رمز یکبار مصرف ارسال شد",
		"userExists": resp,
	})
}

func (ac *AuthController) VerifyOTP(c *gin.Context) {
	var body serializers_auth.VerifyOTPRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}

	resp, err := ac.authService.ValidateOTP(body.Email, body.OTP)
	if err != nil {
		if err == redis.Nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "کد یکبار مصرف منقضی شده است",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	if resp {
		var exists bool = true
		
		if err := ac.authService.CheckExistance(body.Email); errors.Is(err, gorm.ErrRecordNotFound) {
			exists = false
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "مشکل غیرمنتظره ای رخ داده است",
			})
			log.Println(err.Error())
			return
		}
		if !exists {
			err := ac.authService.CreateUser(body.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "مشکل غیرمنتظره ای رخ داده است",
				})
				log.Println(err.Error())
				return
			}
		}

		user, _ := ac.authService.GetUserByGmail(body.Email)
		// TODO: Move logic to service
		
		if jwt, err := ac.authService.GenerateJWT(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "مشکل غیرمنتظره ای رخ داده است",
			})
			log.Println(err.Error())
			return
		}else{
			c.SetCookie("pashmak_authentication", jwt, int(ac.AppConfig.TokenAge), "/", ac.AppConfig.CookieDomain, true, false)
			c.SetSameSite(http.SameSiteNoneMode)
		}
		
		
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "ورود با موفقیت انجام شد.",
		})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "رمز یکبار مصرف اشتباه وارد شده.",
		})
		return
	}
}

func (ac *AuthController) ProtectedRouter(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "این یک api محافظت شده است :)",
	})
}

func (ac *AuthController) LoginWithPassword(c *gin.Context) {
	var body serializers_auth.LoginWithPasswordRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}

	jwt, err := ac.authService.LoginWithPassword(body.Email, body.Password)
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" || err.Error() == "user has no password" || err.Error() == "record not found" { // TODO: Integrate errors
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "نام کاربری یا رمز عبور اشتباه است",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.SetCookie("pashmak_authentication", jwt, int(ac.AppConfig.TokenAge), "/", ac.AppConfig.CookieDomain, true, false)
	c.SetSameSite(http.SameSiteNoneMode)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "ورود با موفقیت انجام شد.",
	})
}

func (ac *AuthController) ForgetPassword(c *gin.Context) {
	var body serializers_auth.ForgetPasswordRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}

	err := ac.authService.ForgetPassword(body.Email)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "کاربر پیدا نشد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "خطای غیر منتظره رخ داد.",
		})
		log.Println(err.Error())
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "کد تایید به ایمیل ارسال شد.",
		})
	}
}

func (ac *AuthController) ForgetPasswordVerify(c *gin.Context) {
	var body serializers_auth.ForgetPasswordVerifyRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}

	jwt, resp, err := ac.authService.VerifyForgetPassword(body.Email, body.OTP)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "کاربر پیدا نشد",
			})
			return
		}
		if err.Error() == "OTP expired" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "کد تایید منقضی شده است",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	if resp {
		c.SetCookie("pashmak_authentication", jwt, int(ac.AppConfig.TokenAge), "/", ac.AppConfig.CookieDomain, true, false)
		c.SetSameSite(http.SameSiteNoneMode)
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "رمز یکبار مصرف صحیح وارد شده.",
		})
		return
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "رمز یکبار مصرف اشتباه وارد شده.",
		})
		return
	}
}

func (ac *AuthController) ForgetPasswordReset(c *gin.Context) {
	var body serializers_auth.ForgetPasswordResetRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
		})
		return
	}
	value, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "شمامجاز به انجام این عملیات نمی باشید.",
		})
		return
	}
	userinfo := value.(services_auth.UserInfo)
	err := ac.authService.ResetForgetPassword(userinfo, body.Password)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "کاربر پیدا نشد",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیرمنتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "رمز عبور با موفقیت تغییر یافت",
		})
	}
}

func (ac *AuthController) SignUp(c *gin.Context) {
	var body serializers_auth.SignUpRequest
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "در خواندن بدنه ی درخواست خطایی رخ داد",
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
	userpayload := userinfo.(services_auth.UserInfo)
	err := ac.authService.SignUp(userpayload, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "مشکل غیر منتظره ای رخ داده است",
		})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "ثبت نام با موفقیت انجام شد.",
	})
}
	