package repository

import (
	"context"
	"github.com/watchlist-kata/auth/internal/models"
)

type AuthRepository interface {
	AddRefreshToken(ctx context.Context, userId, email, rt string) error
	СheckRefreshToken(ctx context.Context, rt string) (models.User, error)
}
