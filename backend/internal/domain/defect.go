package domain

import "time"

// DefectListItem describes a subset of fields for table/list views.
type DefectListItem struct {
	ID           string
	ProjectID    string
	ProjectName  string
	Title        string
	Priority     string
	Status       string
	AssigneeID   string
	AssigneeName string
	DueDate      *time.Time
	UpdatedAt    time.Time
}

// Defect represents full defect entity with description/comments info.
type Defect struct {
	ID          string
	ProjectID   string
	ProjectName string
	Title       string
	Description string
	Priority    string
	Severity    string
	Status      string
	AssigneeID  string
	Assignee    string
	DueDate     *time.Time
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Attachments []Attachment
	Comments    []Comment
}

// DefectCreate describes payload for creating a defect.
type DefectCreate struct {
	ProjectID   string
	Title       string
	Description string
	Priority    string
	Severity    string
	AssigneeID  string
	DueDate     *time.Time
	CreatedBy   string
}

// DefectFilter keeps optional filters passed from transport layer.
type DefectFilter struct {
	Status   string
	Priority string
	Limit    int
	Project  string
}

type Comment struct {
	ID         string
	AuthorID   string
	AuthorName string
	Body       string
	CreatedAt  time.Time
}

type CommentCreate struct {
	DefectID string
	AuthorID string
	Body     string
}

type Attachment struct {
	DefectID    string
	ID          string
	Filename    string
	ContentType string
	SizeBytes   int64
	StorageKey  string
	UploadedAt  time.Time
}

type AttachmentCreate struct {
	DefectID    string
	Filename    string
	ContentType string
	SizeBytes   int64
	StorageKey  string
}
