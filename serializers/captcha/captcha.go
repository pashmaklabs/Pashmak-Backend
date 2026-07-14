package captcha

type VerificationRequest struct {
    Token string `json:"token" binding:"required"`
}

type VerificationResponse struct {
    Success bool     `json:"success"`
    Message string   `json:"message,omitempty"`
    Errors  []string `json:"errors,omitempty"`
}