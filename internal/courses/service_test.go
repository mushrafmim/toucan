package courses

import (
	"context"
	"testing"
	"toucan/internal/enrollments"
	"toucan/internal/identity"
	"toucan/internal/users"
)

type mockUserService struct {
	users.Service
}

func (m *mockUserService) GetByExternalSubject(ctx context.Context, subject string) (users.User, error) {
	return users.User{
		ID:              "user-123",
		ExternalSubject: subject,
		Roles:           []users.Role{users.RoleInstructor},
	}, nil
}

func TestServiceCreateAndPublishCourse(t *testing.T) {
	ctx := identity.ContextWithPrincipal(context.Background(), identity.Principal{
		Subject: "auth0|123",
		Roles:   []string{"admin"},
	})
	userRepo := &mockUserService{}
	enrollmentRepo := enrollments.NewMemoryRepository()
	enrollmentService := enrollments.NewService(enrollmentRepo, userRepo)
	courseRepo := NewMemoryRepository()
	service := NewService(courseRepo, userRepo, enrollmentService)

	course, err := service.Create(ctx, CreateCourseInput{
		Title:       " Distributed Systems Basics ",
		Summary:     "Core distributed systems concepts.",
		Description: "A course on replicas, consistency, and failure handling.",
		Category:    "Architecture",
		Tags:        []string{"Distributed", "architecture", "distributed"},
	})
	if err != nil {
		t.Fatalf("create course: %v", err)
	}

	if course.Status != StatusDraft {
		t.Fatalf("expected draft status, got %q", course.Status)
	}
	if course.Slug != "distributed-systems-basics" {
		t.Fatalf("expected slug to be normalized, got %q", course.Slug)
	}
	if len(course.Tags) != 2 {
		t.Fatalf("expected unique normalized tags, got %v", course.Tags)
	}

	published, err := service.Publish(ctx, course.ID)
	if err != nil {
		t.Fatalf("publish course: %v", err)
	}
	if published.Status != StatusPublished {
		t.Fatalf("expected published status, got %q", published.Status)
	}
	if published.PublishedAt.IsZero() {
		t.Fatal("expected publish timestamp to be set")
	}
}

func TestServiceRejectsInvalidCourse(t *testing.T) {
	ctx := identity.ContextWithPrincipal(context.Background(), identity.Principal{
		Subject: "auth0|123",
		Roles:   []string{"admin"},
	})
	userRepo := &mockUserService{}
	enrollmentRepo := enrollments.NewMemoryRepository()
	enrollmentService := enrollments.NewService(enrollmentRepo, userRepo)
	courseRepo := NewMemoryRepository()
	service := NewService(courseRepo, userRepo, enrollmentService)

	_, err := service.Create(ctx, CreateCourseInput{Title: "No Summary"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
