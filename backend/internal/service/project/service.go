package project

import (
	"context"
	"strings"

	"defect-tracker/internal/domain"
)

type Repository interface {
	List(ctx context.Context) ([]domain.Project, error)
	Create(ctx context.Context, payload domain.ProjectCreate) (domain.Project, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]domain.Project, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, payload domain.ProjectCreate) (domain.Project, error) {
	if strings.TrimSpace(payload.Name) == "" {
		return domain.Project{}, ErrValidation
	}
	if payload.Stage == "" {
		payload.Stage = "Не указан"
	}
	return s.repo.Create(ctx, payload)
}

var ErrValidation = domainError("validation error")

type domainError string

func (e domainError) Error() string {
	return string(e)
}
