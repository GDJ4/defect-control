package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"defect-tracker/internal/domain"
)

var ErrTokenExpired = errors.New("refresh token expired")

type Repository interface {
	Save(ctx context.Context, userID, token string, expiresAt time.Time) error
	Get(ctx context.Context, token string) (domain.RefreshToken, error)
	Delete(ctx context.Context, token string) error
	DeleteByUser(ctx context.Context, userID string) error
}

type Service struct {
	repo Repository
	ttl  time.Duration
}

func NewService(repo Repository, ttl time.Duration) *Service {
	return &Service{repo: repo, ttl: ttl}
}

func (s *Service) Issue(ctx context.Context, userID string) (domain.RefreshToken, error) {
	token := generateToken()
	expiresAt := time.Now().Add(s.ttl)

	if err := s.repo.Save(ctx, userID, token, expiresAt); err != nil {
		return domain.RefreshToken{}, err
	}

	return domain.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *Service) Rotate(ctx context.Context, token string) (domain.RefreshToken, error) {
	existing, err := s.repo.Get(ctx, token)
	if err != nil {
		return domain.RefreshToken{}, err
	}
	if time.Now().After(existing.ExpiresAt) {
		return domain.RefreshToken{}, ErrTokenExpired
	}

	if err := s.repo.Delete(ctx, token); err != nil {
		return domain.RefreshToken{}, err
	}

	return s.Issue(ctx, existing.UserID)
}

func (s *Service) Revoke(ctx context.Context, token string) error {
	return s.repo.Delete(ctx, token)
}

func (s *Service) RevokeUserTokens(ctx context.Context, userID string) error {
	return s.repo.DeleteByUser(ctx, userID)
}

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
