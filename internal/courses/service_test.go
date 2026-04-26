package courses

import "testing"

func TestServiceCreateAndPublishCourse(t *testing.T) {
	service := NewService(NewMemoryRepository())

	course, err := service.Create(CreateCourseInput{
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

	published, err := service.Publish(course.ID)
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
	service := NewService(NewMemoryRepository())

	_, err := service.Create(CreateCourseInput{Title: "No Summary"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
