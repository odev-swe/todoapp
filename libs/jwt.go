package libs

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Data any `json:"data,omitempty"`
}

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

func GenerateToken(data any, secret string, tokenType TokenType) (string, error) {
	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Data:             data,
	}

	switch tokenType {
	case AccessToken:
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(15 * time.Second))
	case RefreshToken:
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)

	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
