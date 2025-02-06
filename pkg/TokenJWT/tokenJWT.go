package TokenJWT

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenJWT struct {
	secretKey string
	ttl       time.Duration
}

func NewTokenJWT(secretKey string, ttl time.Duration) *TokenJWT {
	return &TokenJWT{secretKey: secretKey, ttl: ttl}
}

func (t *TokenJWT) GenerateToken(payload map[string]string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(t.ttl).Unix(),
	}

	for key, val := range payload {
		claims[key] = val
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.secretKey))
}

func (t *TokenJWT) ValidateJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		return false, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return false, fmt.Errorf("token has expired")
			}
		}
		return true, nil
	}

	return false, fmt.Errorf("token is invalid")
}
