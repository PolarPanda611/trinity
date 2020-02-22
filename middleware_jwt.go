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
	expireTime := GetCurrentTime().Add(time.Duration(GlobalTrinity.setting.Security.Authentication.JwtExpireHour) * time.Hour)

	claims := Claims{
		userkey,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    GlobalTrinity.setting.Security.Authentication.JwtIssuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, rErr := tokenClaims.SignedString([]byte(GlobalTrinity.setting.Security.Authentication.SecretKey))
	if rErr != nil {
		return "", rErr, ErrGeneratedToken
	}
	return token, nil, nil
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GlobalTrinity.setting.Security.Authentication.SecretKey), nil
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
	if prefix != GlobalTrinity.setting.Security.Authentication.JwtHeaderPrefix {
		return nil, ErrTokenWrongHeaderPrefix, ErrUnverifiedToken
	}
	tokenClaims, rErr, uErr := ParseToken(token)
	if rErr != nil {
		return nil, rErr, uErr
	}
	if !tokenClaims.StandardClaims.VerifyIssuer(GlobalTrinity.setting.Security.Authentication.JwtIssuer, true) {
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
func ParseUnverifiedToken(token string) (*FedidClaims, error) {
	p := new(jwt.Parser)
	p.SkipClaimsValidation = true
	claim := FedidClaims{}
	_, _, err := p.ParseUnverified(token, &claim)
	if err != nil {
		return nil, err
	}
	if GlobalTrinity.setting.Security.Authentication.JwtVerifyExpireHour {
		if !claim.StandardClaims.VerifyExpiresAt(GetCurrentTimeUnix(), true) {
			return nil, ErrTokenExpired
		}
	}
	if GlobalTrinity.setting.Security.Authentication.JwtVerifyIssuer {
		if !claim.StandardClaims.VerifyIssuer(GlobalTrinity.setting.Security.Authentication.JwtIssuer, true) {
			return nil, ErrTokenWrongIssuer
		}
	}

	return &claim, nil

}

// CheckUnverifiedTokenValid check authorization header token is valid
func CheckUnverifiedTokenValid(c *gin.Context) (*FedidClaims, error) {
	if c.Request.Header.Get("Authorization") == "" || len(strings.Fields(c.Request.Header.Get("Authorization"))) != 2 {
		return nil, ErrTokenWrongAuthorization
	}
	prefix := strings.Fields(c.Request.Header.Get("Authorization"))[0]
	token := strings.Fields(c.Request.Header.Get("Authorization"))[1]
	if prefix != GlobalTrinity.setting.Security.Authentication.JwtHeaderPrefix {
		return nil, ErrTokenWrongHeaderPrefix
	}
	tokenClaims, err := ParseUnverifiedToken(token)
	if err != nil {
		return nil, err
	}
	if !tokenClaims.StandardClaims.VerifyIssuer(GlobalTrinity.setting.Security.Authentication.JwtIssuer, true) {
		return nil, ErrTokenWrongIssuer
	}
	return tokenClaims, nil
}

// JwtUnverifiedAuthBackend get claim info
func JwtUnverifiedAuthBackend(c *gin.Context) error {
	tokenClaims, err := CheckUnverifiedTokenValid(c)
	if err != nil {
		return err
	}
	c.Set("Username", tokenClaims.UID)
	return nil

}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, rErr, _ := CheckTokenValid(c)

		if rErr != nil {
			c.AbortWithError(401, rErr)
			return
		}

		c.Next()
	}
}
