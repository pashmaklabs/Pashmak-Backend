package serializers_auth

type StartEmailAuthRequest struct {
	Email string `json:"email"`
}

type StartEmailAuthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}
