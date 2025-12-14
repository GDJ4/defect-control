package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"

	"defect-tracker/internal/domain"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) List(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, stage, description, start_date, end_date, created_by, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var (
			project domain.Project
			start   sql.NullTime
			end     sql.NullTime
		)

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Stage,
			&project.Description,
			&start,
			&end,
			&project.CreatedBy,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if start.Valid {
			project.StartDate = &start.Time
		}
		if end.Valid {
			project.EndDate = &end.Time
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) Create(ctx context.Context, payload domain.ProjectCreate) (domain.Project, error) {
	var project domain.Project

	err := r.pool.QueryRow(ctx, `
		INSERT INTO projects (name, stage, description, start_date, end_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		payload.Name,
		payload.Stage,
		payload.Description,
		dateOrNil(payload.StartDate),
		dateOrNil(payload.EndDate),
		payload.CreatedBy,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return domain.Project{}, err
	}

	project.Name = payload.Name
	project.Stage = payload.Stage
	project.Description = payload.Description
	project.StartDate = payload.StartDate
	project.EndDate = payload.EndDate
	project.CreatedBy = payload.CreatedBy

	return project, nil
}
