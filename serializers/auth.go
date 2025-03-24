package serializers_auth

type SendOTPRequest struct {
	Email string `json:"email"`
}

type SendOTPResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type VerifyOTPRequest struct{
	Email string	`json:"email"`
	OTP string		`json:"otp"`
}

type ForgetPasswordRequest struct{
	Email string	`json:"email"`
}