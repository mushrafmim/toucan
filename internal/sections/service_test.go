package sections

import (
	"context"
	"testing"

	"toucan/internal/courses"
	"toucan/internal/enrollments"
	"toucan/internal/identity"
	"toucan/internal/users"
)

type mockUserServiceInServiceTest struct {
	users.Service
}

func (m *mockUserServiceInServiceTest) GetByExternalSubject(ctx context.Context, subject string) (users.User, error) {
	return users.User{
		ID:              "user-123",
		ExternalSubject: subject,
		Roles:           []users.Role{users.RoleInstructor},
	}, nil
}

func TestServiceManagesSections(t *testing.T) {
	ctx := identity.ContextWithPrincipal(context.Background(), identity.Principal{
		Subject: "auth0|123",
		Roles:   []string{"admin"},
	})
	userRepo := &mockUserServiceInServiceTest{}
	enrollmentRepo := enrollments.NewMemoryRepository()
	enrollmentService := enrollments.NewService(enrollmentRepo, userRepo)
	courseService := courses.NewService(courses.NewMemoryRepository(), userRepo, enrollmentService)
	service := NewService(NewMemoryRepository(), courseService)

	course, err := courseService.Create(ctx, courses.CreateCourseInput{
		Title:       "Platform Onboarding",
		Summary:     "Everything a new learner needs first.",
		Description: "Covers setup, navigation, and delivery expectations.",
	})
	if err != nil {
		t.Fatalf("create course: %v", err)
	}

	section, err := service.Create(ctx, CreateSectionInput{
		CourseID: course.ID,
		Title:    "Week 1",
		Summary:  "Orientation materials.",
		Position: 1,
	})
	if err != nil {
		t.Fatalf("create section: %v", err)
	}
	if section.Title != "Week 1" {
		t.Fatalf("expected section title, got %q", section.Title)
	}
	storedSection, err := service.Get(ctx, section.ID)
	if err != nil {
		t.Fatalf("get section: %v", err)
	}
	if storedSection.CourseID != course.ID {
		t.Fatalf("expected course id %q, got %q", course.ID, storedSection.CourseID)
	}
}
