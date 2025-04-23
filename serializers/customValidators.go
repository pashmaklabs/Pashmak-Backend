package serializers

import "github.com/go-playground/validator/v10"

func RegisterCustomValidators(v *validator.Validate) {
    v.RegisterValidation("containsuppercase", func(fl validator.FieldLevel) bool {
        password := fl.Field().String()
        for _, char := range password {
            if char >= 'A' && char <= 'Z' {
                return true
            }
        }
        return false
    })

    v.RegisterValidation("containsnumber", func(fl validator.FieldLevel) bool {
        password := fl.Field().String()
        for _, char := range password {
            if char >= '0' && char <= '9' {
                return true
            }
        }
        return false
    })
}