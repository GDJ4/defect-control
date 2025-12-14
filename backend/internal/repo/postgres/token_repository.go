package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"defect-tracker/internal/domain"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

func (r *TokenRepository) Save(ctx context.Context, userID, token string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *TokenRepository) Get(ctx context.Context, token string) (domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, token, expires_at, created_at
		FROM auth_tokens WHERE token = $1`,
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	return rt, err
}

func (r *TokenRepository) Delete(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM auth_tokens WHERE token = $1`, token)
	return err
}

func (r *TokenRepository) DeleteByUser(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM auth_tokens WHERE user_id = $1`, userID)
	return err
}
