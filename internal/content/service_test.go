package content

import (
	"testing"

	"toucan/internal/courses"
	"toucan/internal/sections"
)

func TestServiceManagesContentItems(t *testing.T) {
	courseService := courses.NewService(courses.NewMemoryRepository())
	sectionService := sections.NewService(sections.NewMemoryRepository(), courseService)
	service := NewService(NewMemoryRepository(), sectionService)

	course, err := courseService.Create(courses.CreateCourseInput{
		Title:       "Platform Onboarding",
		Summary:     "Everything a new learner needs first.",
		Description: "Covers setup, navigation, and delivery expectations.",
	})
	if err != nil {
		t.Fatalf("create course: %v", err)
	}
	section, err := sectionService.Create(sections.CreateSectionInput{
		CourseID: course.ID,
		Title:    "Week 1",
		Summary:  "Orientation materials.",
		Position: 1,
	})
	if err != nil {
		t.Fatalf("create section: %v", err)
	}

	item, err := service.Create(CreateItemInput{
		SectionID: section.ID,
		Title:     "Welcome Video",
		Summary:   "Overview of the platform.",
		Type:      TypeVideo,
		Position:  1,
		SourceURL: "https://cdn.example.test/welcome.mp4",
		Metadata:  map[string]any{"duration_seconds": 180},
	})
	if err != nil {
		t.Fatalf("create content item: %v", err)
	}
	if item.Type != TypeVideo {
		t.Fatalf("expected content type video, got %q", item.Type)
	}
	if item.SectionID != section.ID {
		t.Fatalf("expected section id %q, got %q", section.ID, item.SectionID)
	}
}
