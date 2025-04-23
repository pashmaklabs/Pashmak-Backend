package middlewares_validation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var errorMessages = map[string]string{
	"Password.password_complexity": "Password must be at least 8 characters long and contain at least one uppercase letter, one number, and one special character!",
	"Email.email":                  "Email has wrong format!",
	"OTP.min":                      "OTP is a 4 digit number!",
	"OTP.numeric":                  "OTP is a 4 digit number!",
}

// ValidationMiddleware validates the request body against the provided struct type
func ValidationMiddleware[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		// Bind JSON payload
		if err := c.ShouldBindJSON(&req); err != nil {
			// Handle binding errors (e.g., malformed JSON)
			if _, ok := err.(validator.ValidationErrors); !ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "در خواندن بدنه ی درخواست مشکلی پیش آمده است",
					"errors":  map[string]string{"payload": err.Error()},
				})
				c.Abort()
				return
			}

			// Handle validation errors
			errors := make(map[string]string)
			for _, e := range err.(validator.ValidationErrors) {
				key := e.Field() + "." + e.Tag()
				if msg, exists := errorMessages[key]; exists {
					errors[e.Field()] = msg
				} else {
					errors[e.Field()] = e.Error() // Fallback to default validator message
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

		// Store validated struct in context
		c.Set("validated", req)
		c.Next()
	}
}
