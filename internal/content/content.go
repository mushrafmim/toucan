package content

import (
	"errors"
	"strings"
	"time"

	"toucan/internal/shared"
)

type Type string

const (
	TypeVideo    Type = "video"
	TypePDF      Type = "pdf"
	TypeRichText Type = "rich_text"
	TypeFile     Type = "file"
	TypeLink     Type = "link"
	TypeEmbed    Type = "embed"
)

type Item struct {
	ID        string         `json:"id"`
	SectionID string         `json:"section_id"`
	Title     string         `json:"title"`
	Summary   string         `json:"summary"`
	Type      Type           `json:"type"`
	Position  int            `json:"position"`
	Configs   map[string]any `json:"configs,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type CreateItemInput struct {
	SectionID string         `json:"section_id"`
	Title     string         `json:"title"`
	Summary   string         `json:"summary"`
	Type      Type           `json:"type"`
	Position  int            `json:"position"`
	Configs   map[string]any `json:"configs"`
}

type UpdateItemInput struct {
	Title    string         `json:"title"`
	Summary  string         `json:"summary"`
	Type     Type           `json:"type"`
	Position int            `json:"position"`
	Configs  map[string]any `json:"configs"`
}

type ListFilter struct {
	SectionID string
	Page      int
	PageSize  int
}

type ListResult = shared.ListResult[Item]

var (
	ErrNotFound    = errors.New("content item not found")
	ErrInvalidType = errors.New("invalid content item type")
	ErrValidation  = errors.New("validation failed")
)

func cloneConfigs(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	clone := make(map[string]any, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}

func validateType(contentType Type) error {
	switch contentType {
	case TypeVideo, TypePDF, TypeRichText, TypeFile, TypeLink, TypeEmbed:
		return nil
	default:
		return ErrInvalidType
	}
}
