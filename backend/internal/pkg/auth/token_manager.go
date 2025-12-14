package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"defect-tracker/internal/domain"
)

type Claims struct {
	Role     string `json:"role"`
	FullName string `json:"fullName"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *Manager) Generate(user domain.User) (string, error) {
	now := time.Now()
	claims := Claims{
		Role:     user.Role,
		FullName: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(token string) (Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return *claims, nil
	}
	return Claims{}, jwt.ErrTokenInvalidClaims
}
