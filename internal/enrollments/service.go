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

func (s *Service) Delete(ctx context.Context, courseID, userID string) error {
	return s.repo.Delete(ctx, courseID, userID)
}

func (s *Service) ListByCourse(ctx context.Context, courseID string) ([]Enrollment, error) {
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

func (s *Service) DB() *sql.DB {
	return s.repo.DB()
}

func (s *Service) WithTx(tx *sql.Tx) *Service {
	return &Service{
		repo:     s.repo.WithTx(tx),
		userRepo: s.userRepo,
	}
}
