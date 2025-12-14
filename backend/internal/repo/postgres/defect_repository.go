package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"defect-tracker/internal/domain"
)

// DefectRepository works with defects table/views.
type DefectRepository struct {
	pool *pgxpool.Pool
}

func NewDefectRepository(pool *pgxpool.Pool) *DefectRepository {
	return &DefectRepository{pool: pool}
}

func (r *DefectRepository) List(ctx context.Context, filter domain.DefectFilter) ([]domain.DefectListItem, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT
			d.id,
			d.project_id,
			COALESCE(p.name, '') AS project_name,
			d.title,
			d.priority,
			d.status,
			COALESCE(d.assignee_id::text, '') AS assignee_id,
			COALESCE(u.full_name, '') AS assignee_name,
			d.due_date,
			d.updated_at
		FROM defects d
		LEFT JOIN projects p ON p.id = d.project_id
		LEFT JOIN users u ON u.id = d.assignee_id
		WHERE 1=1
	`)

	args := make([]any, 0, 4)
	argPos := 1

	if filter.Project != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND d.project_id = $%d", argPos))
		args = append(args, filter.Project)
		argPos++
	}

	if filter.Status != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND d.status = $%d", argPos))
		args = append(args, filter.Status)
		argPos++
	}

	if filter.Priority != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND d.priority = $%d", argPos))
		args = append(args, filter.Priority)
		argPos++
	}

	queryBuilder.WriteString(" ORDER BY d.created_at DESC")
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", argPos))
	args = append(args, filter.Limit)

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.DefectListItem, 0)
	for rows.Next() {
		var (
			item     domain.DefectListItem
			dueDate  sql.NullTime
			assignee sql.NullString
		)
		if err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.Title,
			&item.Priority,
			&item.Status,
			&assignee,
			&item.AssigneeName,
			&dueDate,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if dueDate.Valid {
			item.DueDate = &dueDate.Time
		}
		if assignee.Valid {
			item.AssigneeID = assignee.String
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *DefectRepository) Create(ctx context.Context, payload domain.DefectCreate) (domain.Defect, error) {
	var (
		defect domain.Defect
		due    sql.NullTime
	)
	if payload.DueDate != nil {
		due.Valid = true
		due.Time = *payload.DueDate
	}

	err := r.pool.QueryRow(ctx, `
		INSERT INTO defects (project_id, title, description, priority, severity, status, assignee_id, due_date, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, 'NEW', $6, $7, $8, $8)
		RETURNING id, created_at, updated_at`,
		payload.ProjectID,
		payload.Title,
		payload.Description,
		payload.Priority,
		payload.Severity,
		nullIfEmpty(payload.AssigneeID),
		dateOrNil(payload.DueDate),
		payload.CreatedBy,
	).Scan(&defect.ID, &defect.CreatedAt, &defect.UpdatedAt)
	if err != nil {
		return domain.Defect{}, err
	}

	defect.ProjectID = payload.ProjectID
	defect.Title = payload.Title
	defect.Description = payload.Description
	defect.Priority = payload.Priority
	defect.Severity = payload.Severity
	defect.Status = "NEW"
	defect.AssigneeID = payload.AssigneeID
	defect.DueDate = payload.DueDate
	defect.CreatedBy = payload.CreatedBy

	return defect, nil
}

func (r *DefectRepository) GetByID(ctx context.Context, id string) (domain.Defect, error) {
	var (
		defect  domain.Defect
		due     sql.NullTime
		assID   sql.NullString
		assName sql.NullString
	)

	err := r.pool.QueryRow(ctx, `
		SELECT
			d.id,
			d.project_id,
			COALESCE(p.name, '') AS project_name,
			d.title,
			d.description,
			d.priority,
			d.severity,
			d.status,
			d.assignee_id,
			COALESCE(u.full_name, '') AS assignee_name,
			d.due_date,
			d.created_by,
			d.created_at,
			d.updated_at
		FROM defects d
		LEFT JOIN projects p ON p.id = d.project_id
		LEFT JOIN users u ON u.id = d.assignee_id
		WHERE d.id = $1`,
		id,
	).Scan(
		&defect.ID,
		&defect.ProjectID,
		&defect.ProjectName,
		&defect.Title,
		&defect.Description,
		&defect.Priority,
		&defect.Severity,
		&defect.Status,
		&assID,
		&assName,
		&due,
		&defect.CreatedBy,
		&defect.CreatedAt,
		&defect.UpdatedAt,
	)
	if err != nil {
		return domain.Defect{}, err
	}

	if due.Valid {
		defect.DueDate = &due.Time
	}
	if assID.Valid {
		defect.AssigneeID = assID.String
	}
	if assName.Valid {
		defect.Assignee = assName.String
	}

	attachments, err := r.ListAttachments(ctx, id)
	if err != nil {
		return domain.Defect{}, err
	}
	comments, err := r.ListComments(ctx, id)
	if err != nil {
		return domain.Defect{}, err
	}
	defect.Attachments = attachments
	defect.Comments = comments

	return defect, nil
}

func (r *DefectRepository) AddComment(ctx context.Context, payload domain.CommentCreate) (domain.Comment, error) {
	var comment domain.Comment
	err := r.pool.QueryRow(ctx, `
		INSERT INTO defect_comments (defect_id, author_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		payload.DefectID,
		payload.AuthorID,
		payload.Body,
	).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return domain.Comment{}, err
	}

	comment.AuthorID = payload.AuthorID
	comment.Body = payload.Body

	err = r.pool.QueryRow(ctx, `SELECT full_name FROM users WHERE id = $1`, payload.AuthorID).Scan(&comment.AuthorName)
	if err != nil {
		comment.AuthorName = "Неизвестно"
	}
	return comment, nil
}

func (r *DefectRepository) ListComments(ctx context.Context, defectID string) ([]domain.Comment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.author_id, COALESCE(u.full_name, ''), c.body, c.created_at
		FROM defect_comments c
		LEFT JOIN users u ON u.id = c.author_id
		WHERE c.defect_id = $1
		ORDER BY c.created_at ASC`,
		defectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.AuthorName, &comment.Body, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (r *DefectRepository) AddAttachment(ctx context.Context, payload domain.AttachmentCreate) (domain.Attachment, error) {
	var attachment domain.Attachment
	err := r.pool.QueryRow(ctx, `
		INSERT INTO defect_attachments (defect_id, filename, content_type, size_bytes, storage_key)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		payload.DefectID,
		payload.Filename,
		payload.ContentType,
		payload.SizeBytes,
		payload.StorageKey,
	).Scan(&attachment.ID, &attachment.UploadedAt)
	if err != nil {
		return domain.Attachment{}, err
	}

	attachment.DefectID = payload.DefectID
	attachment.Filename = payload.Filename
	attachment.ContentType = payload.ContentType
	attachment.SizeBytes = payload.SizeBytes
	attachment.StorageKey = payload.StorageKey
	return attachment, nil
}

func (r *DefectRepository) ListAttachments(ctx context.Context, defectID string) ([]domain.Attachment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT defect_id, id, filename, content_type, size_bytes, storage_key, created_at
		FROM defect_attachments
		WHERE defect_id = $1
		ORDER BY created_at DESC`,
		defectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []domain.Attachment
	for rows.Next() {
		var att domain.Attachment
		if err := rows.Scan(&att.DefectID, &att.ID, &att.Filename, &att.ContentType, &att.SizeBytes, &att.StorageKey, &att.UploadedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, att)
	}
	return attachments, rows.Err()
}

func (r *DefectRepository) GetAttachment(ctx context.Context, defectID, attachmentID string) (domain.Attachment, error) {
	var att domain.Attachment
	err := r.pool.QueryRow(ctx, `
		SELECT defect_id, id, filename, content_type, size_bytes, storage_key, created_at
		FROM defect_attachments
		WHERE defect_id = $1 AND id = $2`,
		defectID, attachmentID,
	).Scan(&att.DefectID, &att.ID, &att.Filename, &att.ContentType, &att.SizeBytes, &att.StorageKey, &att.UploadedAt)
	return att, err
}

func (r *DefectRepository) UpdateStatus(ctx context.Context, id, status, actorID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE defects SET status = $1, updated_by = $2, updated_at = NOW()
		WHERE id = $3`,
		status, actorID, id,
	)
	return err
}

func (r *DefectRepository) AddHistory(ctx context.Context, defectID, actorID, field, oldValue, newValue string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO defect_history (defect_id, actor_id, field, old_value, new_value)
		VALUES ($1, $2, $3, $4, $5)`,
		defectID, actorID, field, oldValue, newValue,
	)
	return err
}
func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func dateOrNil(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t
}
