package sections

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"toucan/internal/courses"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleListSections(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(ListFilter{
		CourseID: r.URL.Query().Get("course_id"),
		Page:     parseInt(r.URL.Query().Get("page"), 1),
		PageSize: parseInt(r.URL.Query().Get("page_size"), 10),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) HandleGetSection(w http.ResponseWriter, r *http.Request) {
	section, err := h.service.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, section)
}

func (h *Handler) HandleCreateSection(w http.ResponseWriter, r *http.Request) {
	var input CreateSectionInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}

	section, err := h.service.Create(input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, section)
}

func (h *Handler) HandleUpdateSection(w http.ResponseWriter, r *http.Request) {
	var input UpdateSectionInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}

	section, err := h.service.Update(r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, section)
}

func (h *Handler) HandleDeleteSection(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
	case errors.Is(err, ErrNotFound), errors.Is(err, courses.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrValidation):
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
