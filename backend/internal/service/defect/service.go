package defect

import (
	"context"
	"fmt"
	"strings"

	"defect-tracker/internal/domain"
)

type Repository interface {
	List(ctx context.Context, filter domain.DefectFilter) ([]domain.DefectListItem, error)
	Create(ctx context.Context, payload domain.DefectCreate) (domain.Defect, error)
	GetByID(ctx context.Context, id string) (domain.Defect, error)
	UpdateStatus(ctx context.Context, id, status, actorID string) error
	AddHistory(ctx context.Context, defectID, actorID, field, oldValue, newValue string) error
	AddComment(ctx context.Context, payload domain.CommentCreate) (domain.Comment, error)
	ListComments(ctx context.Context, defectID string) ([]domain.Comment, error)
	AddAttachment(ctx context.Context, payload domain.AttachmentCreate) (domain.Attachment, error)
	ListAttachments(ctx context.Context, defectID string) ([]domain.Attachment, error)
	GetAttachment(ctx context.Context, defectID, attachmentID string) (domain.Attachment, error)
}

// Service handles domain-level logic for defects.
type Service struct {
	repo Repository
}

var (
	allowedStatuses = map[string]struct{}{
		"NEW":         {},
		"IN_PROGRESS": {},
		"IN_REVIEW":   {},
		"CLOSED":      {},
		"CANCELED":    {},
	}
	allowedPriorities = map[string]struct{}{
		"LOW":      {},
		"MEDIUM":   {},
		"HIGH":     {},
		"CRITICAL": {},
	}
)

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filter domain.DefectFilter) ([]domain.DefectListItem, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}

	filter.Status = normalizeEnum(filter.Status, allowedStatuses)
	filter.Priority = normalizeEnum(filter.Priority, allowedPriorities)

	return s.repo.List(ctx, filter)
}

func (s *Service) Create(ctx context.Context, payload domain.DefectCreate) (domain.Defect, error) {
	payload.Priority = normalizeEnum(payload.Priority, allowedPriorities)
	payload.Severity = normalizeEnum(payload.Severity, allowedPriorities)
	return s.repo.Create(ctx, payload)
}

func (s *Service) Get(ctx context.Context, id string) (domain.Defect, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) AddComment(ctx context.Context, payload domain.CommentCreate) (domain.Comment, error) {
	if strings.TrimSpace(payload.Body) == "" {
		return domain.Comment{}, fmt.Errorf("comment body is empty")
	}
	return s.repo.AddComment(ctx, payload)
}

func (s *Service) ListComments(ctx context.Context, defectID string) ([]domain.Comment, error) {
	return s.repo.ListComments(ctx, defectID)
}

func (s *Service) AddAttachment(ctx context.Context, payload domain.AttachmentCreate) (domain.Attachment, error) {
	if payload.SizeBytes <= 0 {
		return domain.Attachment{}, fmt.Errorf("attachment is empty")
	}
	return s.repo.AddAttachment(ctx, payload)
}

func (s *Service) ListAttachments(ctx context.Context, defectID string) ([]domain.Attachment, error) {
	return s.repo.ListAttachments(ctx, defectID)
}

func (s *Service) GetAttachment(ctx context.Context, defectID, attachmentID string) (domain.Attachment, error) {
	return s.repo.GetAttachment(ctx, defectID, attachmentID)
}

func (s *Service) UpdateStatus(ctx context.Context, defectID string, actor domain.User, newStatus string) (domain.Defect, error) {
	nextStatus := normalizeEnum(newStatus, allowedStatuses)
	if nextStatus == "" {
		return domain.Defect{}, fmt.Errorf("некорректный статус")
	}

	defect, err := s.repo.GetByID(ctx, defectID)
	if err != nil {
		return domain.Defect{}, err
	}

	if defect.Status == nextStatus {
		return defect, nil
	}

	if !canTransition(defect.Status, nextStatus) {
		return domain.Defect{}, fmt.Errorf("переход %s -> %s запрещён", defect.Status, nextStatus)
	}

	if requiresManager(nextStatus) && actor.Role != "manager" {
		return domain.Defect{}, fmt.Errorf("статус %s доступен только менеджеру", nextStatus)
	}

	if err := s.repo.UpdateStatus(ctx, defectID, nextStatus, actor.ID); err != nil {
		return domain.Defect{}, err
	}

	_ = s.repo.AddHistory(ctx, defectID, actor.ID, "status", defect.Status, nextStatus)

	return s.repo.GetByID(ctx, defectID)
}

func canTransition(current, next string) bool {
	transitions := map[string][]string{
		"NEW":         {"IN_PROGRESS", "CANCELED"},
		"IN_PROGRESS": {"IN_REVIEW", "CANCELED"},
		"IN_REVIEW":   {"CLOSED", "CANCELED"},
		"CLOSED":      {},
		"CANCELED":    {},
	}
	for _, allowed := range transitions[current] {
		if allowed == next {
			return true
		}
	}
	return false
}

func requiresManager(status string) bool {
	return status == "CLOSED" || status == "CANCELED"
}

func normalizeEnum(value string, allowed map[string]struct{}) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	if _, ok := allowed[value]; ok {
		return value
	}
	return ""
}
