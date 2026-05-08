package courses

import (
	"errors"
	"strings"
	"time"

	"toucan/internal/shared"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Level string

const (
	LevelBeginner     Level = "beginner"
	LevelIntermediate Level = "intermediate"
	LevelAdvanced     Level = "advanced"
)

type Course struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Level       Level     `json:"level"`
	Tags        []string  `json:"tags"`
	Status      Status    `json:"status"`
	CreatorID   string    `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

type CreateCourseInput struct {
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Level       Level    `json:"level"`
	Tags        []string `json:"tags"`
}

type UpdateCourseInput struct {
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Level       Level    `json:"level"`
	Tags        []string `json:"tags"`
}

type ListFilter struct {
	Query    string
	Status   Status
	Page     int
	PageSize int
}

type ListResult = shared.ListResult[Course]

var (
	ErrNotFound          = errors.New("course not found")
	ErrInvalidStatus     = errors.New("invalid course status")
	ErrInvalidLevel      = errors.New("invalid course level")
	ErrValidation        = errors.New("validation failed")
	ErrInvalidTransition = errors.New("invalid status transition")
)

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(tags))
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		normalized = append(normalized, tag)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

func validateLevel(level Level) error {
	switch level {
	case "", LevelBeginner, LevelIntermediate, LevelAdvanced:
		return nil
	default:
		return ErrInvalidLevel
	}
}

func validateStatus(status Status) error {
	switch status {
	case "", StatusDraft, StatusPublished, StatusArchived:
		return nil
	default:
		return ErrInvalidStatus
	}
}
