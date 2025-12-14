package domain

import (
	"errors"
	"time"
)

// User represents an application user.
type User struct {
	ID           string
	Email        string
	FullName     string
	Role         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// UserRegister describes payload for public registration.
type UserRegister struct {
	Email    string
	FullName string
	Password string
	Role     string
}

// ErrEmailAlreadyExists indicates unique constraint violation for user email.
var ErrEmailAlreadyExists = errors.New("user with email already exists")
