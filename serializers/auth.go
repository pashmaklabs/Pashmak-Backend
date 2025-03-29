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
	FirstName       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}
