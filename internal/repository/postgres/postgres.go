package postgres

import (
	"context"
	"database/sql"
	"github.com/watchlist-kata/auth/internal/models"
)

type PostgesRepositoryImpl struct {
	db *sql.DB
}

func NewPostgesRepositoryImpl(db *sql.DB) *PostgesRepositoryImpl {
	return &PostgesRepositoryImpl{db: db}
}

func (p *PostgesRepositoryImpl) AddRefreshToken(ctx context.Context, userId, email, rt string) error {
	query := `INSERT INTO refresh_token(user_id, email, refresh) VALUES ($1, $2, $3)`
	_, err := p.db.ExecContext(ctx, query, userId, email, rt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgesRepositoryImpl) СheckRefreshToken(ctx context.Context, rt string) (models.User, error) {
	var user models.User
	query := `SELECT user_id, email FROM refresh_token WHERE refresh = $1`
	err := p.db.QueryRowContext(ctx, query, rt).Scan(&user.UserId, &user.Email)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
