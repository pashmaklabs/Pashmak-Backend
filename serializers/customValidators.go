package serializers

import (
	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
)

var PasswordRegex = regexp2.MustCompile(`^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$`, regexp2.None)

func RegisterCustomValidators(v *validator.Validate) {
    v.RegisterValidation("password_complexity", func(fl validator.FieldLevel) bool {
        password := fl.Field().String()

        if match, _ := PasswordRegex.MatchString(password); !match {
            return false
        }
        return true
    })
}