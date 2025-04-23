package middlewares_validation

import (
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
                "message": "در خواندن بدنه ی درخواست خطایی رخ داد",
            })
            c.Abort()
            return
        }


        // Store the validated struct in the context
        c.Set("validated", value)
        c.Next()
		
    }
}