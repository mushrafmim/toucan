package sections

import (
	"testing"

	"toucan/internal/courses"
)

func TestServiceManagesSections(t *testing.T) {
	courseService := courses.NewService(courses.NewMemoryRepository())
	service := NewService(NewMemoryRepository(), courseService)

	course, err := courseService.Create(courses.CreateCourseInput{
		Title:       "Platform Onboarding",
		Summary:     "Everything a new learner needs first.",
		Description: "Covers setup, navigation, and delivery expectations.",
	})
	if err != nil {
		t.Fatalf("create course: %v", err)
	}

	section, err := service.Create(CreateSectionInput{
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
	storedSection, err := service.Get(section.ID)
	if err != nil {
		t.Fatalf("get section: %v", err)
	}
	if storedSection.CourseID != course.ID {
		t.Fatalf("expected course id %q, got %q", course.ID, storedSection.CourseID)
	}
}
