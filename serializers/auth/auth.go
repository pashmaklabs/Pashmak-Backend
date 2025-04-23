package serializers_authtype

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,min=4,max=4"`
}

type LoginWithPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,containsany=!@#$%^&*,containsuppercase,containsnumber"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgetPasswordVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,min=4,max=4"`
}

type ForgetPasswordResetRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,containsany=!@#$%^&*,containsuppercase,containsnumber"`
}

type SignUpRequest struct {
	FirstName       string `json:"firstname" binding:"required,min=2,max=50"`
	LastName        string `json:"lastname" binding:"required,min=2,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,containsany=!@#$%^&*,containsuppercase,containsnumber"`
	PasswordConfirm string `json:"password_confirm" binding:"required,min=8,containsany=!@#$%^&*,containsuppercase,containsnumber"`
}
