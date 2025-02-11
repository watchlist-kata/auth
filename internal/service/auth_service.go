package service

import (
	"context"
	"github.com/watchlist-kata/auth/internal/repository"
	"github.com/watchlist-kata/auth/pkg/tokenJWT"
)

type AuthService interface {
	GenerateTokens(ctx context.Context, userId, email string) (string, string, error)
	RefreshToken(ctx context.Context, rt string) (string, error)
	ValidateJWT(token string) error
}

type AuthSeviceImpl struct {
	TokenJWT *tokenJWT.TokenJWT
	repo     repository.AuthRepository
}

func NewAuthSeviceImpl(tokenJWT *tokenJWT.TokenJWT, repo repository.AuthRepository) *AuthSeviceImpl {
	return &AuthSeviceImpl{TokenJWT: tokenJWT, repo: repo}
}

func (a *AuthSeviceImpl) GenerateTokens(ctx context.Context, userId, email string) (string, string, error) {
	refreshToken := a.TokenJWT.GenerateRefreshToken()
	err := a.repo.AddRefreshToken(ctx, userId, email, refreshToken)
	if err != nil {
		return "", "", err
	}

	jwt, err := a.TokenJWT.GenerateToken(map[string]string{"email": email, "userId": userId})
	if err != nil {
		return "", "", err
	}

	return jwt, refreshToken, nil
}

func (a *AuthSeviceImpl) RefreshToken(ctx context.Context, rt string) (string, error) {
	user, err := a.repo.СheckRefreshToken(ctx, rt)
	if err != nil {
		return "", err
	}

	jwt, err := a.TokenJWT.GenerateToken(map[string]string{"email": user.Email, "userId": user.UserId})
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (a *AuthSeviceImpl) ValidateJWT(token string) error {
	_, err := a.TokenJWT.ValidateJWT(token)
	if err != nil {
		return err
	}

	return nil
}
