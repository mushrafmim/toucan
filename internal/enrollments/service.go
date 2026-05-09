package enrollments

import (
	"context"
	"database/sql"
	"fmt"
	"toucan/internal/identity"
	"toucan/internal/users"
)

type Service struct {
	repo     Repository
	userRepo users.Service
}

func NewService(repo Repository, userRepo users.Service) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *Service) Get(ctx context.Context, courseID, userID string) (Enrollment, error) {
	return s.repo.Get(ctx, courseID, userID)
}

func (s *Service) Create(ctx context.Context, enrollment Enrollment) (Enrollment, error) {
	return s.repo.Create(ctx, enrollment)
}

func (s *Service) CreateWithAuth(ctx context.Context, enrollment Enrollment) (Enrollment, error) {
	if err := s.authorize(ctx, enrollment.CourseID, RoleOwner); err != nil {
		return Enrollment{}, err
	}
	return s.repo.Create(ctx, enrollment)
}

func (s *Service) Delete(ctx context.Context, courseID, userID string) error {
	return s.repo.Delete(ctx, courseID, userID)
}

func (s *Service) DeleteWithAuth(ctx context.Context, courseID, userID string) error {
	if err := s.authorize(ctx, courseID, RoleOwner); err != nil {
		return err
	}
	return s.repo.Delete(ctx, courseID, userID)
}

func (s *Service) ListByCourse(ctx context.Context, courseID string) ([]Enrollment, error) {
	return s.repo.ListByCourse(ctx, courseID)
}

func (s *Service) ListByCourseWithAuth(ctx context.Context, courseID string) ([]Enrollment, error) {
	if err := s.authorize(ctx, courseID, RoleManager); err != nil {
		return nil, err
	}
	return s.repo.ListByCourse(ctx, courseID)
}

func (s *Service) GetMyEnrollment(ctx context.Context, courseID string) (Enrollment, error) {
	principal, ok := identity.PrincipalFromContext(ctx)
	if !ok {
		return Enrollment{}, fmt.Errorf("unauthorized")
	}

	user, err := s.userRepo.GetByExternalSubject(ctx, principal.Subject)
	if err != nil {
		return Enrollment{}, fmt.Errorf("resolve user: %w", err)
	}

	return s.repo.Get(ctx, courseID, user.ID)
}

func (s *Service) authorize(ctx context.Context, courseID string, requiredRole Role) error {
	principal, ok := identity.PrincipalFromContext(ctx)
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	// Global admins can do anything
	for _, role := range principal.Roles {
		if role == "admin" {
			return nil
		}
	}

	user, err := s.userRepo.GetByExternalSubject(ctx, principal.Subject)
	if err != nil {
		return fmt.Errorf("unauthorized")
	}

	enrollment, err := s.repo.Get(ctx, courseID, user.ID)
	if err != nil {
		return fmt.Errorf("unauthorized")
	}

	// Owner can do anything a manager can do
	if enrollment.Role == RoleOwner {
		return nil
	}

	if requiredRole == RoleManager && enrollment.Role == RoleManager {
		return nil
	}

	return fmt.Errorf("unauthorized")
}

func (s *Service) DB() *sql.DB {
	return s.repo.DB()
}

func (s *Service) WithTx(tx *sql.Tx) *Service {
	return &Service{
		repo:     s.repo.WithTx(tx),
		userRepo: s.userRepo,
	}
}
