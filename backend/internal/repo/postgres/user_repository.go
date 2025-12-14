package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"defect-tracker/internal/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, full_name, role, password_hash, created_at, updated_at
		FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, full_name, role, password_hash, created_at, updated_at
		FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET password_hash = $1, updated_at = NOW()
		WHERE id = $2`,
		passwordHash, id,
	)
	return err
}

func (r *UserRepository) Create(ctx context.Context, email, fullName, role, passwordHash string) (domain.User, error) {
	var user domain.User
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (email, full_name, role, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, full_name, role, password_hash, created_at, updated_at`,
		email, fullName, role, passwordHash,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrEmailAlreadyExists
		}
		return domain.User{}, err
	}
	return user, nil
}
