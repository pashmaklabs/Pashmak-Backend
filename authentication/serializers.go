package authentication

type StartEmailAuthRequest struct{
	Email string	`json:"email"`
}

type StartEmailAuthResponse struct{
	//
}

type VerifyOTPRequest struct{
	OTP string	`json:"OTP"`
}