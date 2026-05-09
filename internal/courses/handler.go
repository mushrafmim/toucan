package courses

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"toucan/internal/enrollments"
)

type Handler struct {
	service           *Service
	enrollmentService *enrollments.Service
	logger            *log.Logger
}

func NewHandler(service *Service, enrollmentService *enrollments.Service, logger *log.Logger) *Handler {
	return &Handler{service: service, enrollmentService: enrollmentService, logger: logger}
}

func (h *Handler) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) HandleRoot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"name":        "Toucan",
		"description": "Generic LMS backend",
	})
}

func (h *Handler) HandleListCourses(w http.ResponseWriter, r *http.Request) {
	filter := ListFilter{
		Query:    r.URL.Query().Get("q"),
		Status:   Status(r.URL.Query().Get("status")),
		Page:     parseInt(r.URL.Query().Get("page"), 1),
		PageSize: parseInt(r.URL.Query().Get("page_size"), 10),
	}

	result, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) HandleListMyCourses(w http.ResponseWriter, r *http.Request) {
	filter := ListFilter{
		Query:    r.URL.Query().Get("q"),
		Status:   Status(r.URL.Query().Get("status")),
		Page:     parseInt(r.URL.Query().Get("page"), 1),
		PageSize: parseInt(r.URL.Query().Get("page_size"), 10),
	}

	result, err := h.service.ListMyCourses(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) HandleGetCourse(w http.ResponseWriter, r *http.Request) {
	course, err := h.service.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, course)
}

func (h *Handler) HandleGetMyMembership(w http.ResponseWriter, r *http.Request) {
	member, err := h.enrollmentService.GetMyEnrollment(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, member)
}

func (h *Handler) HandleCreateCourse(w http.ResponseWriter, r *http.Request) {
	var input CreateCourseInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}

	course, err := h.service.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, course)
}

func (h *Handler) HandleUpdateCourse(w http.ResponseWriter, r *http.Request) {
	var input UpdateCourseInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}

	course, err := h.service.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, course)
}

func (h *Handler) HandleDeleteCourse(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandlePublishCourse(w http.ResponseWriter, r *http.Request) {
	course, err := h.service.Publish(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, course)
}

func (h *Handler) HandleArchiveCourse(w http.ResponseWriter, r *http.Request) {
	course, err := h.service.Archive(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, course)
}

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func decodeJSON(r *http.Request, out any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

func parseInt(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, enrollments.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		status = http.StatusForbidden
	case errors.Is(err, ErrValidation), errors.Is(err, ErrInvalidStatus), errors.Is(err, ErrInvalidLevel), errors.Is(err, ErrInvalidTransition):
		status = http.StatusBadRequest
	case errors.As(err, new(*json.SyntaxError)), errors.As(err, new(*json.UnmarshalTypeError)):
		status = http.StatusBadRequest
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}
