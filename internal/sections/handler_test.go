package sections

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"toucan/internal/courses"
)

func TestHandlerSectionLifecycle(t *testing.T) {
	courseService := courses.NewService(courses.NewMemoryRepository())
	courseHandler := courses.NewHandler(courseService, log.New(io.Discard, "", 0))
	sectionHandler := NewHandler(NewService(NewMemoryRepository(), courseService))

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
	createCourseReq := httptest.NewRequest(http.MethodPost, "/api/v1/courses", bytes.NewReader(courseBody))
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

	var section Section
	if err := json.Unmarshal(createSectionRes.Body.Bytes(), &section); err != nil {
		t.Fatalf("decode section: %v", err)
	}

	listSectionsReq := httptest.NewRequest(http.MethodGet, "/api/v1/sections?course_id="+createdCourse.ID, nil)
	listSectionsRes := httptest.NewRecorder()
	handler.ServeHTTP(listSectionsRes, listSectionsReq)
	if listSectionsRes.Code != http.StatusOK {
		t.Fatalf("expected 200 from list sections, got %d", listSectionsRes.Code)
	}

	getSectionReq := httptest.NewRequest(http.MethodGet, "/api/v1/sections/"+section.ID, nil)
	getSectionRes := httptest.NewRecorder()
	handler.ServeHTTP(getSectionRes, getSectionReq)
	if getSectionRes.Code != http.StatusOK {
		t.Fatalf("expected 200 from get section, got %d", getSectionRes.Code)
	}

	var hydrated Section
	if err := json.Unmarshal(getSectionRes.Body.Bytes(), &hydrated); err != nil {
		t.Fatalf("decode hydrated section: %v", err)
	}
	if hydrated.CourseID != createdCourse.ID {
		t.Fatalf("expected section course id %q, got %q", createdCourse.ID, hydrated.CourseID)
	}
}
