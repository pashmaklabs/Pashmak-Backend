package middlewares_validation

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationMiddleware(v interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := reflect.New(reflect.TypeOf(v)).Interface()

		if err := c.ShouldBindJSON(value); err != nil {
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				errors := make(map[string]string)
				for _, e := range validationErrs {
					errors[e.Field()] = e.Tag()
					if e.Field() == "Password" && e.Tag() == "password_complexity" {
                        errors[e.Field()] = "Password must be at least 8 characters long and contain at least one uppercase letter, one number, and one special character!"
                    }
					if e.Field() == "Email" && e.Tag() == "email" {
                        errors[e.Field()] = "Email has wrong format!"
                    }
					if e.Field() == "OTP" {
						if e.Tag() == "min" || e.Tag() == "numeric"{
							errors[e.Field()] = "OTP is a 4 digit number!"
						}
                        
                    }
				}
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "داده‌های ورودی نامعتبر است",
					"errors":  errors,
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "مشکل غیر منتظره ای پیش آمده است",
			})
			log.Println("Error in ValidationMiddleware!")
			c.Abort()
			return
		}

		// Store the validated struct in the context
		c.Set("validated", value)
		c.Next()

	}
}
