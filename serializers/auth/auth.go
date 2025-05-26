package serializers_authtype

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,min=4,max=4,numeric"`
}

type LoginWithPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password_complexity"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgetPasswordVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,min=4,max=4,numeric"`
}

type ForgetPasswordResetRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password_complexity"`
}

type SignUpRequest struct {
	FirstName       string `json:"firstname" binding:"min=0,max=50"`
	LastName        string `json:"lastname" binding:"min=0,max=50"`
	Password        string `json:"password" binding:"required,password_complexity"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,password_complexity"`
}
