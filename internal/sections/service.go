package sections

import (
	"context"
	"fmt"
	"strings"
	"time"

	"toucan/internal/courses"
)

type CourseLookup interface {
	Get(ctx context.Context, id string) (courses.Course, error)
}

type Service struct {
	repo         *Repository
	courseLookup CourseLookup
}

func NewService(repo *Repository, courseLookup CourseLookup) *Service {
	return &Service{repo: repo, courseLookup: courseLookup}
}

func (s *Service) List(ctx context.Context, filter ListFilter) (ListResult, error) {
	if filter.CourseID != "" {
		if _, err := s.courseLookup.Get(ctx, strings.TrimSpace(filter.CourseID)); err != nil {
			return ListResult{}, err
		}
	}
	return s.repo.List(filter), nil
}

func (s *Service) Get(ctx context.Context, id string) (Section, error) {
	return s.repo.Get(strings.TrimSpace(id))
}

func (s *Service) Create(ctx context.Context, input CreateSectionInput) (Section, error) {
	if err := validateSectionInput(input); err != nil {
		return Section{}, err
	}
	courseID := strings.TrimSpace(input.CourseID)
	if _, err := s.courseLookup.Get(ctx, courseID); err != nil {
		return Section{}, err
	}

	now := time.Now().UTC()
	section := Section{
		CourseID:  courseID,
		Title:     normalizeText(input.Title),
		Summary:   normalizeText(input.Summary),
		Position:  normalizePosition(input.Position, 1),
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(section)
}

func (s *Service) Update(ctx context.Context, id string, input UpdateSectionInput) (Section, error) {
	if err := validateSectionUpdateInput(input); err != nil {
		return Section{}, err
	}

	section, err := s.repo.Get(strings.TrimSpace(id))
	if err != nil {
		return Section{}, err
	}
	section.Title = normalizeText(input.Title)
	section.Summary = normalizeText(input.Summary)
	section.Position = normalizePosition(input.Position, section.Position)
	return s.repo.Update(section)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(strings.TrimSpace(id))
}

func validateSectionInput(input CreateSectionInput) error {
	if strings.TrimSpace(input.CourseID) == "" || strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: course_id and title are required", ErrValidation)
	}
	if input.Position < 0 {
		return fmt.Errorf("%w: section position must be zero or greater", ErrValidation)
	}
	return nil
}

func validateSectionUpdateInput(input UpdateSectionInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrValidation)
	}
	if input.Position < 0 {
		return fmt.Errorf("%w: section position must be zero or greater", ErrValidation)
	}
	return nil
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
