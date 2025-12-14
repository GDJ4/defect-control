package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"defect-tracker/internal/domain"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	Create(ctx context.Context, email, fullName, role, passwordHash string) (domain.User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (domain.User, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	if cleanEmail == "" || strings.TrimSpace(password) == "" {
		return domain.User{}, ErrInvalidCredentials
	}

	user, err := s.repo.GetByEmail(ctx, cleanEmail)
	if err != nil {
		return domain.User{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, ErrInvalidCredentials
	}
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Register(ctx context.Context, payload domain.UserRegister) (domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	fullName := strings.TrimSpace(payload.FullName)
	password := strings.TrimSpace(payload.Password)
	role := strings.TrimSpace(strings.ToLower(payload.Role))

	if email == "" || fullName == "" || password == "" {
		return domain.User{}, fmt.Errorf("email, ФИО и пароль обязательны")
	}

	if len(password) < 6 {
		return domain.User{}, fmt.Errorf("пароль должен быть не короче 6 символов")
	}

	if role == "" {
		role = "engineer"
	}
	if !isRoleAllowed(role) {
		return domain.User{}, fmt.Errorf("недопустимая роль")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user, err := s.repo.Create(ctx, email, fullName, role, string(hash))
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return domain.User{}, fmt.Errorf("пользователь с таким email уже зарегистрирован")
		}
		return domain.User{}, err
	}
	return user, nil
}

func (s *Service) ChangePassword(ctx context.Context, id, currentPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, string(newHash))
}

func isRoleAllowed(role string) bool {
	switch role {
	case "manager", "engineer", "observer":
		return true
	default:
		return false
	}
}
