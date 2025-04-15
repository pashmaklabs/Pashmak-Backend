package services_auth

import (
	"errors"
)

var ErrAuth = struct {
	ErrNoPassword       error
	ErrOTPExpired       error
	ErrOperationTimeOut error
	ErrPasswordMismatch error
	ErrTokenExpired error
	ErrJWTParseFailure error
}{
	ErrNoPassword:       errors.New("user has no password"),
	ErrOTPExpired:       errors.New("OTP expired"),
	ErrOperationTimeOut: errors.New("operation timed out"),
	ErrPasswordMismatch: errors.New("passwords do not match"),
	ErrTokenExpired: errors.New("Token is expired"),
	ErrJWTParseFailure: errors.New("Couldn't parse claims"),
}
