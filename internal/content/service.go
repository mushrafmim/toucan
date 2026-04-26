package content

import (
	"fmt"
	"strings"
	"time"

	"toucan/internal/sections"
)

type SectionLookup interface {
	Get(id string) (sections.Section, error)
}

type Service struct {
	repo          *Repository
	sectionLookup SectionLookup
}

func NewService(repo *Repository, sectionLookup SectionLookup) *Service {
	return &Service{repo: repo, sectionLookup: sectionLookup}
}

func (s *Service) List(filter ListFilter) (ListResult, error) {
	if filter.SectionID != "" {
		if _, err := s.sectionLookup.Get(strings.TrimSpace(filter.SectionID)); err != nil {
			return ListResult{}, err
		}
	}
	return s.repo.List(filter), nil
}

func (s *Service) Get(id string) (Item, error) {
	return s.repo.Get(strings.TrimSpace(id))
}

func (s *Service) Create(input CreateItemInput) (Item, error) {
	if err := validateCreateInput(input); err != nil {
		return Item{}, err
	}
	sectionID := strings.TrimSpace(input.SectionID)
	if _, err := s.sectionLookup.Get(sectionID); err != nil {
		return Item{}, err
	}

	now := time.Now().UTC()
	item := Item{
		SectionID: sectionID,
		Title:     normalizeText(input.Title),
		Summary:   normalizeText(input.Summary),
		Type:      input.Type,
		Position:  normalizePosition(input.Position, 1),
		SourceURL: normalizeText(input.SourceURL),
		Body:      normalizeText(input.Body),
		Metadata:  cloneMetadata(input.Metadata),
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(item)
}

func (s *Service) Update(id string, input UpdateItemInput) (Item, error) {
	if err := validateUpdateInput(input); err != nil {
		return Item{}, err
	}
	item, err := s.repo.Get(strings.TrimSpace(id))
	if err != nil {
		return Item{}, err
	}
	item.Title = normalizeText(input.Title)
	item.Summary = normalizeText(input.Summary)
	item.Type = input.Type
	item.Position = normalizePosition(input.Position, item.Position)
	item.SourceURL = normalizeText(input.SourceURL)
	item.Body = normalizeText(input.Body)
	item.Metadata = cloneMetadata(input.Metadata)
	return s.repo.Update(item)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(strings.TrimSpace(id))
}

func validateCreateInput(input CreateItemInput) error {
	if strings.TrimSpace(input.SectionID) == "" || strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: section_id and title are required", ErrValidation)
	}
	if input.Position < 0 {
		return fmt.Errorf("%w: content item position must be zero or greater", ErrValidation)
	}
	return validateType(input.Type)
}

func validateUpdateInput(input UpdateItemInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrValidation)
	}
	if input.Position < 0 {
		return fmt.Errorf("%w: content item position must be zero or greater", ErrValidation)
	}
	return validateType(input.Type)
}

func normalizePosition(position, fallback int) int {
	if position > 0 {
		return position
	}
	if fallback > 0 {
		return fallback
	}
	return 1
}
