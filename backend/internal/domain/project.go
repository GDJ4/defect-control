package domain

import "time"

type Project struct {
	ID          string
	Name        string
	Stage       string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProjectCreate struct {
	Name        string
	Stage       string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedBy   string
}
