package services_auth

import (
	models_auth "pashmak.com/pashmak/models"
	"github.com/golang-jwt/jwt"
	"pashmak.com/pashmak/bootstrap"
	"crypto/x509"
    "encoding/pem"
	"time"
	"crypto/rsa"
	"log"
	"os"
	"github.com/google/uuid"
	"errors"
)

var signKey *rsa.PrivateKey

type UserInfo struct {
	Email 	string
	ID 		uint
}

type CustomClaim struct {
	*jwt.StandardClaims
	*UserInfo
}

func LoadPrivateKey() {
    privateKeyFile, err := os.Open(bootstrap.PRIVATE_KEY_PATH)
    if err != nil {
        log.Fatalf("Error opening private key file: %v", err)
    }
    defer privateKeyFile.Close()

    pemFileInfo, _ := privateKeyFile.Stat()
    var size = pemFileInfo.Size()
    pemBytes := make([]byte, size)

    _, err = privateKeyFile.Read(pemBytes)
    if err != nil {
        log.Fatalf("Error reading private key file: %v", err)
    }

    data, _ := pem.Decode(pemBytes)
    privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
    if err != nil {
        log.Fatalf("Error parsing private key: %v", err)
    }

    signKey = privateKeyImported
}

func (as *AuthService)GenerateJWT(user models_auth.User) (string, error){
	Id := uuid.New().String()
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	
	t.Claims = &CustomClaim{
		&jwt.StandardClaims{
			Id: Id, ExpiresAt: time.Now().Add(time.Second * time.Duration(as.AppConfig.TokenAge)).Unix(),
		},
		&UserInfo{
			Email: user.Email,
			ID: user.ID,
		},
	}
	
	if signKey == nil {
		LoadPrivateKey()
	}
	return t.SignedString(signKey)
}

func (as *AuthService) ParseToken(tokenString string) (*CustomClaim, error) {
	if signKey == nil {	
		LoadPrivateKey()
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return signKey.Public(), nil
	})
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*CustomClaim)
	if !ok {	
		return nil, errors.New("Couldn't parse claims")
	}
	return claim, nil
}

func (as *AuthService)VerifyJWT(tokenString string) (*CustomClaim, error) {
    claim, err := as.ParseToken(tokenString)
    if err != nil {
        return nil, err
    }

    if claim.StandardClaims.ExpiresAt < time.Now().Unix() {
        return nil, errors.New("Token is expired")
    }


	// [FIXME] : this part is not working so i pass it to you hossein :)
    // _, err = as.GetJWTBlacklistByJTI(claim.StandardClaims.Id)
    // if err != nil { 
		// }
		// return nil, err
	return claim, nil
    
}

func (as *AuthService) GetJWTBlacklistByJTI(jti string) (models_auth.JWTBlacklist, error) {
    var jwtBlacklist models_auth.JWTBlacklist
    result := as.DB.First(&jwtBlacklist, "jti = ?", jti)
    return jwtBlacklist, result.Error
}
