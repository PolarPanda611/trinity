package trinity

import (
	"errors"
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
func GenerateToken(userkey string) (string, error) {
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
	token, err := tokenClaims.SignedString([]byte(DefaultSecretkey))
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(DefaultSecretkey), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// CheckTokenValid check authorization header token is valid
func CheckTokenValid(c *gin.Context) (*Claims, error) {
	if c.Request.Header.Get("Authorization") == "" || len(strings.Fields(c.Request.Header.Get("Authorization"))) != 2 {
		return nil, errors.New("app.err.loaddatafailed")
	}
	prefix := strings.Fields(c.Request.Header.Get("Authorization"))[0]
	token := strings.Fields(c.Request.Header.Get("Authorization"))[1]
	if prefix != DefaultJwtheaderprefix {
		return nil, errors.New("app.err.loaddatafailed")
	}
	tokenClaims, err := ParseToken(token)
	if err != nil {
		return nil, err
	}
	if !tokenClaims.StandardClaims.VerifyIssuer(DefaultJwtissuer, true) {
		return nil, errors.New("app.err.failedtoverifytoken")
	}
	return tokenClaims, nil
}

// JwtAuthBackend check authorization header token is valid
func JwtAuthBackend(c *gin.Context) error {
	tokenClaims, err := CheckTokenValid(c)
	if err != nil {
		return errors.New("app.err.failedtoverifytoken")
	}
	c.Set("UserID", tokenClaims.Userkey)
	return nil

}

// ParseUnverifiedToken parsing token
func ParseUnverifiedToken(token string) (*FedidClaims, error) {
	p := new(jwt.Parser)
	p.SkipClaimsValidation = true
	claim := FedidClaims{}
	_, _, err := p.ParseUnverified(token, &claim)
	if err != nil {
		return nil, err
	}
	if !claim.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("app.err.failedtoverifytoken")
	}
	if !claim.StandardClaims.VerifyIssuer(DefaultJwtissuer, true) {
		return nil, errors.New("app.err.failedtoverifytoken")
	}
	return &claim, nil

}

// CheckUnverifiedTokenValid check authorization header token is valid
func CheckUnverifiedTokenValid(c *gin.Context) (*FedidClaims, error) {
	if c.Request.Header.Get("Authorization") == "" || len(strings.Fields(c.Request.Header.Get("Authorization"))) != 2 {
		return nil, errors.New("app.err.loaddatafailed")
	}
	prefix := strings.Fields(c.Request.Header.Get("Authorization"))[0]
	token := strings.Fields(c.Request.Header.Get("Authorization"))[1]
	if prefix != DefaultJwtheaderprefix {
		return nil, errors.New("app.err.loaddatafailed")
	}
	tokenClaims, err := ParseUnverifiedToken(token)
	if err != nil {
		return nil, err
	}
	if !tokenClaims.StandardClaims.VerifyIssuer(DefaultJwtissuer, true) {
		return nil, errors.New("app.err.failedtoverifytoken")
	}
	return tokenClaims, nil
}

// JwtUnverifiedAuthBackend get claim info
func JwtUnverifiedAuthBackend(c *gin.Context) error {
	tokenClaims, err := CheckUnverifiedTokenValid(c)
	if err != nil {
		return errors.New("app.err.failedtoverifytoken")
	}
	c.Set("UserID", tokenClaims.UID)
	return nil

}
