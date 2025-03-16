package serializers_auth

type SendOTPRequest struct {
	Email string `json:"email"`
}

type SendOTPResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}
