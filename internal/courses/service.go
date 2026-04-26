package courses

import (
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(filter ListFilter) (ListResult, error) {
	if err := validateStatus(filter.Status); err != nil {
		return ListResult{}, err
	}
	return s.repo.List(filter), nil
}

func (s *Service) Get(id string) (Course, error) {
	return s.repo.Get(strings.TrimSpace(id))
}

func (s *Service) Create(input CreateCourseInput) (Course, error) {
	if err := validateCreateInput(input); err != nil {
		return Course{}, err
	}

	now := time.Now().UTC()
	course := Course{
		Title:       normalizeText(input.Title),
		Slug:        slugify(input.Title),
		Summary:     normalizeText(input.Summary),
		Description: normalizeText(input.Description),
		Category:    normalizeText(input.Category),
		Level:       defaultLevel(input.Level),
		Tags:        normalizeTags(input.Tags),
		Status:      StatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.Create(course)
}

func (s *Service) Update(id string, input UpdateCourseInput) (Course, error) {
	if err := validateUpdateInput(input); err != nil {
		return Course{}, err
	}

	course, err := s.repo.Get(strings.TrimSpace(id))
	if err != nil {
		return Course{}, err
	}

	course.Title = normalizeText(input.Title)
	course.Slug = slugify(input.Title)
	course.Summary = normalizeText(input.Summary)
	course.Description = normalizeText(input.Description)
	course.Category = normalizeText(input.Category)
	course.Level = defaultLevel(input.Level)
	course.Tags = normalizeTags(input.Tags)

	return s.repo.Update(course)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(strings.TrimSpace(id))
}

func (s *Service) Publish(id string) (Course, error) {
	course, err := s.repo.Get(strings.TrimSpace(id))
	if err != nil {
		return Course{}, err
	}
	if course.Status == StatusArchived {
		return Course{}, fmt.Errorf("%w: archived course cannot be published", ErrInvalidTransition)
	}
	if course.Status != StatusPublished {
		course.Status = StatusPublished
		course.PublishedAt = time.Now().UTC()
	}
	return s.repo.Update(course)
}

func (s *Service) Archive(id string) (Course, error) {
	course, err := s.repo.Get(strings.TrimSpace(id))
	if err != nil {
		return Course{}, err
	}
	if course.Status == StatusArchived {
		return course, nil
	}
	course.Status = StatusArchived
	return s.repo.Update(course)
}

func validateCreateInput(input CreateCourseInput) error {
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Summary) == "" || strings.TrimSpace(input.Description) == "" {
		return fmt.Errorf("%w: title, summary, and description are required", ErrValidation)
	}
	return validateLevel(input.Level)
}

func validateUpdateInput(input UpdateCourseInput) error {
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Summary) == "" || strings.TrimSpace(input.Description) == "" {
		return fmt.Errorf("%w: title, summary, and description are required", ErrValidation)
	}
	return validateLevel(input.Level)
}

func defaultLevel(level Level) Level {
	if level == "" {
		return LevelBeginner
	}
	return level
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_' || r == '/':
			if b.Len() > 0 && !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "course"
	}
	return slug
}
