package authentication

type StartEmailAuthRequest struct{
	Email string	`json:"email"`
}

type StartEmailAuthResponse struct{
	//
}

type VerifyOTPRequest struct{
	Email string	`json:"email"`
	OTP string		`json:"OTP"`
}