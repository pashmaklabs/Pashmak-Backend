package serializers_auth

type SendOTPRequest struct {
	Email string `json:"email"`
}

type SendOTPResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type LoginWithPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
}

type ForgetPasswordVerifyRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type ForgetPasswordResetRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	FirstName       string `json:"firstname" binding:"required,min=2,max=50"`
    LastName        string `json:"lastname" binding:"required,min=2,max=50"`
    Password        string `json:"password" binding:"required,min=8"`
    PasswordConfirm string `json:"passwordConfirm" binding:"required,min=8"`
}
