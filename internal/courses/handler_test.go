package courses

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerCourseLifecycle(t *testing.T) {
	courseHandler := NewHandler(NewService(NewMemoryRepository()), log.New(io.Discard, "", 0))
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
	handler := courseHandler.LoggingMiddleware(mux)

	createPayload := map[string]any{
		"title":       "API Design",
		"summary":     "How to design maintainable APIs.",
		"description": "Covers REST design, versioning, pagination, and error handling.",
		"category":    "Backend",
		"level":       "intermediate",
		"tags":        []string{"api", "rest"},
	}
	body, _ := json.Marshal(createPayload)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/courses", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	handler.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d with body %s", createRes.Code, createRes.Body.String())
	}

	var created Course
	if err := json.Unmarshal(createRes.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created course: %v", err)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/courses?status=draft", nil)
	listRes := httptest.NewRecorder()
	handler.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("expected 200 from list, got %d", listRes.Code)
	}

	publishReq := httptest.NewRequest(http.MethodPost, "/api/v1/courses/"+created.ID+"/publish", nil)
	publishRes := httptest.NewRecorder()
	handler.ServeHTTP(publishRes, publishReq)
	if publishRes.Code != http.StatusOK {
		t.Fatalf("expected 200 from publish, got %d with body %s", publishRes.Code, publishRes.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/courses/"+created.ID, nil)
	getRes := httptest.NewRecorder()
	handler.ServeHTTP(getRes, getReq)
	if getRes.Code != http.StatusOK {
		t.Fatalf("expected 200 from get, got %d", getRes.Code)
	}
}
