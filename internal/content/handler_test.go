package content

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"toucan/internal/courses"
	"toucan/internal/enrollments"
	"toucan/internal/identity"
	"toucan/internal/sections"
	"toucan/internal/users"
)

type mockUserService struct {
	users.Service
}

func (m *mockUserService) GetByExternalSubject(ctx context.Context, subject string) (users.User, error) {
	return users.User{ID: "user-123"}, nil
}

func TestHandlerContentLifecycle(t *testing.T) {
	userSvc := &mockUserService{}
	enrollmentRepo := enrollments.NewMemoryRepository()
	enrollmentService := enrollments.NewService(enrollmentRepo, userSvc)
	courseRepo := courses.NewMemoryRepository()
	courseService := courses.NewService(courseRepo, userSvc, enrollmentService)
	courseHandler := courses.NewHandler(courseService, enrollmentService, log.New(io.Discard, "", 0))

	sectionRepo := sections.NewMemoryRepository()
	sectionService := sections.NewService(sectionRepo, courseService)
	sectionHandler := sections.NewHandler(sectionService)

	contentRepo := NewMemoryRepository()
	contentHandler := NewHandler(NewService(contentRepo, sectionService))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", courseHandler.HandleHealth)
	mux.HandleFunc("GET /", courseHandler.HandleRoot)
	mux.HandleFunc("GET /api/v1/courses", courseHandler.HandleListCourses)
	mux.HandleFunc("POST /api/v1/courses", courseHandler.HandleCreateCourse)
	mux.HandleFunc("GET /api/v1/courses/{id}", courseHandler.HandleGetCourse)
	mux.HandleFunc("PUT /api/v1/courses/{id}", courseHandler.HandleUpdateCourse)
	mux.HandleFunc("DELETE /api/v1/courses/{id}", courseHandler.HandleDeleteCourse)
	mux.HandleFunc("POST /api/v1/courses/{id}/publish", courseHandler.HandlePublishCourse)
	mux.HandleFunc("POST /api/v1/courses/{id}/archive", courseHandler.HandleArchiveCourse)
	mux.HandleFunc("GET /api/v1/sections", sectionHandler.HandleListSections)
	mux.HandleFunc("POST /api/v1/sections", sectionHandler.HandleCreateSection)
	mux.HandleFunc("GET /api/v1/sections/{id}", sectionHandler.HandleGetSection)
	mux.HandleFunc("PUT /api/v1/sections/{id}", sectionHandler.HandleUpdateSection)
	mux.HandleFunc("DELETE /api/v1/sections/{id}", sectionHandler.HandleDeleteSection)
	mux.HandleFunc("GET /api/v1/content", contentHandler.HandleListContent)
	mux.HandleFunc("POST /api/v1/content", contentHandler.HandleCreateContent)
	mux.HandleFunc("GET /api/v1/content/{id}", contentHandler.HandleGetContent)
	mux.HandleFunc("PUT /api/v1/content/{id}", contentHandler.HandleUpdateContent)
	mux.HandleFunc("DELETE /api/v1/content/{id}", contentHandler.HandleDeleteContent)
	handler := courseHandler.LoggingMiddleware(mux)

	createCoursePayload := map[string]any{
		"title":       "API Design",
		"summary":     "How to design maintainable APIs.",
		"description": "Covers REST design, versioning, pagination, and error handling.",
		"category":    "Backend",
		"level":       "intermediate",
		"tags":        []string{"api", "rest"},
	}
	courseBody, _ := json.Marshal(createCoursePayload)
	ctx := identity.ContextWithPrincipal(context.Background(), identity.Principal{
		Subject: "auth0|123",
		Roles:   []string{"admin"},
	})

	createCourseReq := httptest.NewRequest(http.MethodPost, "/api/v1/courses", bytes.NewReader(courseBody))
	createCourseReq = createCourseReq.WithContext(ctx)
	createCourseReq.Header.Set("Content-Type", "application/json")
	createCourseRes := httptest.NewRecorder()
	handler.ServeHTTP(createCourseRes, createCourseReq)
	if createCourseRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 from create course, got %d with body %s", createCourseRes.Code, createCourseRes.Body.String())
	}

	var createdCourse courses.Course
	if err := json.Unmarshal(createCourseRes.Body.Bytes(), &createdCourse); err != nil {
		t.Fatalf("decode created course: %v", err)
	}

	createSectionPayload := map[string]any{
		"course_id": createdCourse.ID,
		"title":     "Getting Started",
		"summary":   "First section for the course.",
		"position":  1,
	}
	sectionBody, _ := json.Marshal(createSectionPayload)
	createSectionReq := httptest.NewRequest(http.MethodPost, "/api/v1/sections", bytes.NewReader(sectionBody))
	createSectionReq.Header.Set("Content-Type", "application/json")
	createSectionRes := httptest.NewRecorder()
	handler.ServeHTTP(createSectionRes, createSectionReq)
	if createSectionRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 from create section, got %d with body %s", createSectionRes.Code, createSectionRes.Body.String())
	}

	var section sections.Section
	if err := json.Unmarshal(createSectionRes.Body.Bytes(), &section); err != nil {
		t.Fatalf("decode section: %v", err)
	}

	itemPayload := map[string]any{
		"section_id": section.ID,
		"title":      "Course Intro Video",
		"summary":    "A quick welcome video.",
		"type":       "video",
		"position":   1,
		"source_url": "https://cdn.example.test/api-intro.mp4",
		"metadata": map[string]any{
			"duration_seconds": float64(90),
		},
	}
	itemBody, _ := json.Marshal(itemPayload)
	itemReq := httptest.NewRequest(http.MethodPost, "/api/v1/content", bytes.NewReader(itemBody))
	itemReq.Header.Set("Content-Type", "application/json")
	itemRes := httptest.NewRecorder()
	handler.ServeHTTP(itemRes, itemReq)
	if itemRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 from create item, got %d with body %s", itemRes.Code, itemRes.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/content?section_id="+section.ID, nil)
	listRes := httptest.NewRecorder()
	handler.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("expected 200 from list content, got %d", listRes.Code)
	}
}
