package tokenJWT

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrSignatureInvalid = errors.New("invalid signature")
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

func (t *TokenJWT) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrSignatureInvalid) { // Добавлено
			return nil, ErrSignatureInvalid
		}
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func (t *TokenJWT) GenerateRefreshToken() string {
	return uuid.New().String()
}
