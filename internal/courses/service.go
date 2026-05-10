package courses

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"toucan/internal/database"
	"toucan/internal/enrollments"
	"toucan/internal/identity"
	"toucan/internal/users"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Service struct {
	repo              Repository
	userRepo          users.Service
	enrollmentService *enrollments.Service
}

func NewService(repo Repository, userRepo users.Service, enrollmentService *enrollments.Service) *Service {
	return &Service{
		repo:              repo,
		userRepo:          userRepo,
		enrollmentService: enrollmentService,
	}
}

func (s *Service) List(ctx context.Context, filter ListFilter) (ListResult, error) {
	if err := validateStatus(filter.Status); err != nil {
		return ListResult{}, err
	}
	return s.repo.List(filter)
}

func (s *Service) ListMyCourses(ctx context.Context, filter ListFilter) (ListResult, error) {
	principal, ok := identity.PrincipalFromContext(ctx)
	if !ok {
		return ListResult{}, ErrUnauthorized
	}

	user, err := s.userRepo.GetByExternalSubject(ctx, principal.Subject)
	if err != nil {
		return ListResult{}, fmt.Errorf("resolve user: %w", err)
	}

	filter.UserID = user.ID
	return s.repo.List(filter)
}

func (s *Service) Get(ctx context.Context, id string) (Course, error) {
	return s.repo.Get(strings.TrimSpace(id))
}

func (s *Service) Create(ctx context.Context, input CreateCourseInput) (Course, error) {
	principal, ok := identity.PrincipalFromContext(ctx)
	if !ok {
		return Course{}, ErrUnauthorized
	}

	// Only instructors can create courses.
	// Admins are platform managers and should not own courses directly to prevent data discrepancy.
	isInstructor := false
	for _, role := range principal.Roles {
		if role == string(users.RoleInstructor) {
			isInstructor = true
			break
		}
	}
	if !isInstructor {
		return Course{}, ErrUnauthorized
	}

	user, err := s.userRepo.GetByExternalSubject(ctx, principal.Subject)
	if err != nil {
		return Course{}, fmt.Errorf("resolve user: %w", err)
	}

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
		CreatorID:   user.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	var created Course
	err = database.Transact(ctx, s.repo.DB(), func(tx *sql.Tx) error {
		repo := s.repo.WithTx(tx)
		enrollmentRepo := s.enrollmentService.WithTx(tx)

		c, err := repo.Create(course)
		if err != nil {
			return err
		}
		created = c

		// Add creator as owner
		_, err = enrollmentRepo.Create(ctx, enrollments.Enrollment{
			CourseID: created.ID,
			UserID:   user.ID,
			Role:     enrollments.RoleOwner,
		})
		return err
	})

	if err != nil {
		return Course{}, err
	}

	return created, nil
}

func (s *Service) Update(ctx context.Context, id string, input UpdateCourseInput) (Course, error) {
	if err := s.authorize(ctx, id, enrollments.RoleManager); err != nil {
		return Course{}, err
	}

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

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.authorize(ctx, id, enrollments.RoleOwner); err != nil {
		return err
	}
	return s.repo.Delete(strings.TrimSpace(id))
}

func (s *Service) Publish(ctx context.Context, id string) (Course, error) {
	if err := s.authorize(ctx, id, enrollments.RoleManager); err != nil {
		return Course{}, err
	}

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

func (s *Service) Archive(ctx context.Context, id string) (Course, error) {
	if err := s.authorize(ctx, id, enrollments.RoleManager); err != nil {
		return Course{}, err
	}

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

func (s *Service) authorize(ctx context.Context, courseID string, requiredRole enrollments.Role) error {
	principal, ok := identity.PrincipalFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	// Global admins can do anything
	for _, role := range principal.Roles {
		if role == "admin" {
			return nil
		}
	}

	user, err := s.userRepo.GetByExternalSubject(ctx, principal.Subject)
	if err != nil {
		return ErrUnauthorized
	}

	enrollment, err := s.enrollmentService.Get(ctx, courseID, user.ID)
	if err != nil {
		return ErrUnauthorized
	}

	// Owner can do anything a manager can do
	if enrollment.Role == enrollments.RoleOwner {
		return nil
	}

	if requiredRole == enrollments.RoleManager && enrollment.Role == enrollments.RoleManager {
		return nil
	}

	return ErrUnauthorized
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
