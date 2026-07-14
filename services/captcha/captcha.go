package services_captcha

import (
    "errors"
    "log"
    
    "github.com/arcaptcha/arcaptcha-go"
)

type CaptchaService struct {
    client *arcaptcha.Website
}

func NewCaptchaService(siteKey, secretKey string) *CaptchaService {
    return &CaptchaService{
        client: arcaptcha.NewWebsite(siteKey, secretKey),
    }
}

func (cs *CaptchaService) VerifyToken(token string) error {
    if token == "" {
        return errors.New("captcha token is required")
    }
    
    result, err := cs.client.Verify(token)
    if err != nil {
        log.Printf("Error verifying captcha: %v", err)
        return errors.New("captcha verification failed")
    }
    
    if !result.Success {
        log.Printf("Captcha failed with error codes: %v", result.ErrorCodes)
        return errors.New("captcha verification failed")
    }
    
    return nil
}