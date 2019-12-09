package trinity

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims data to sign
type Claims struct {
	Userkey string `json:"userkey"`
	jwt.StandardClaims
}

// FedidClaims data to sign
type FedidClaims struct {
	ClientID string `json:"client_id,omitempty"`
	UID      string `json:"uid,omitempty"`
	Origin   string `json:"origin,omitempty"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(userkey string) (string, error, error) {
	//set expire time
	expireTime := time.Now().Add(time.Duration(DefaultJwtexpirehour) * time.Hour)

	claims := Claims{
		userkey,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    DefaultJwtissuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, rErr := tokenClaims.SignedString([]byte(DefaultSecretkey))
	if rErr != nil {
		return "", rErr, ErrGeneratedToken
	}
	return token, nil, nil
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(DefaultSecretkey), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil, nil
		}
	}

	return nil, err, ErrUnverifiedToken
}

// CheckTokenValid check authorization header token is valid
func CheckTokenValid(c *gin.Context) (*Claims, error, error) {
	if c.Request.Header.Get("Authorization") == "" || len(strings.Fields(c.Request.Header.Get("Authorization"))) != 2 {
		return nil, ErrTokenWrongAuthorization, ErrUnverifiedToken
	}
	prefix := strings.Fields(c.Request.Header.Get("Authorization"))[0]
	token := strings.Fields(c.Request.Header.Get("Authorization"))[1]
	if prefix != DefaultJwtheaderprefix {
		return nil, ErrTokenWrongHeaderPrefix, ErrUnverifiedToken
	}
	tokenClaims, rErr, uErr := ParseToken(token)
	if rErr != nil {
		return nil, rErr, uErr
	}
	if !tokenClaims.StandardClaims.VerifyIssuer(DefaultJwtissuer, true) {
		return nil, ErrTokenWrongIssuer, ErrUnverifiedToken
	}
	return tokenClaims, nil, nil
}

// JwtAuthBackend check authorization header token is valid
func JwtAuthBackend(c *gin.Context) (error, error) {
	tokenClaims, rErr, uErr := CheckTokenValid(c)
	if rErr != nil {
		return rErr, uErr
	}
	c.Set("UserID", tokenClaims.Userkey)
	return nil, nil

}

// ParseUnverifiedToken parsing token
func ParseUnverifiedToken(token string) (*FedidClaims, error, error) {
	p := new(jwt.Parser)
	p.SkipClaimsValidation = true
	claim := FedidClaims{}
	_, _, err := p.ParseUnverified(token, &claim)
	if err != nil {
		return nil, err, ErrUnverifiedToken
	}
	if AppSetting.Security.Authentication.JwtVerifyExpireHour {
		if !claim.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true) {
			return nil, ErrTokenExpired, ErrUnverifiedToken
		}
	}
	if AppSetting.Security.Authentication.JwtVerifyIssuer {
		if !claim.StandardClaims.VerifyIssuer(DefaultJwtissuer, true) {
			return nil, ErrTokenWrongIssuer, ErrUnverifiedToken
		}
	}

	return &claim, nil, nil

}

// CheckUnverifiedTokenValid check authorization header token is valid
func CheckUnverifiedTokenValid(c *gin.Context) (*FedidClaims, error, error) {
	if c.Request.Header.Get("Authorization") == "" || len(strings.Fields(c.Request.Header.Get("Authorization"))) != 2 {
		return nil, ErrTokenWrongAuthorization, ErrUnverifiedToken
	}
	prefix := strings.Fields(c.Request.Header.Get("Authorization"))[0]
	token := strings.Fields(c.Request.Header.Get("Authorization"))[1]
	if prefix != DefaultJwtheaderprefix {
		return nil, ErrTokenWrongHeaderPrefix, ErrUnverifiedToken
	}
	tokenClaims, rErr, uErr := ParseUnverifiedToken(token)
	if rErr != nil {
		return nil, rErr, uErr
	}
	if !tokenClaims.StandardClaims.VerifyIssuer(DefaultJwtissuer, true) {
		return nil, ErrTokenWrongIssuer, ErrUnverifiedToken
	}
	return tokenClaims, nil, nil
}

// JwtUnverifiedAuthBackend get claim info
func JwtUnverifiedAuthBackend(c *gin.Context) (rErr, uErr error) {
	tokenClaims, rErr, uErr := CheckUnverifiedTokenValid(c)
	if rErr != nil {
		return rErr, uErr
	}
	c.Set("UserID", tokenClaims.UID)
	return nil, nil

}
