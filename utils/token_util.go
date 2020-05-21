package utils

import (
	"cicd/config"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenClaims struct {
	ID string `json:"userId"`
	jwt.StandardClaims
}

var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token:")
)

func CreateToken(claims *TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobabConfig.AuthConfig.SigningKey))
}

func ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobabConfig.AuthConfig.SigningKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobabConfig.AuthConfig.SigningKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return CreateToken(claims)
	}
	return "", TokenInvalid
}
