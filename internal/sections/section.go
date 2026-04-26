package sections

import (
	"errors"
	"strings"
	"time"
)

type Section struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSectionInput struct {
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Position int    `json:"position"`
}

type UpdateSectionInput struct {
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Position int    `json:"position"`
}

type ListFilter struct {
	CourseID string
	Page     int
	PageSize int
}

type ListResult struct {
	Items    []Section `json:"items"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Total    int       `json:"total"`
}

var (
	ErrNotFound   = errors.New("section not found")
	ErrValidation = errors.New("validation failed")
)

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}
