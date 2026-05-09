package enrollments

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleCreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var input Enrollment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, err)
		return
	}

	created, err := h.service.CreateWithAuth(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) HandleDeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	userID := r.URL.Query().Get("user_id")

	if courseID == "" || userID == "" {
		h.writeError(w, errors.New("course_id and user_id are required"))
		return
	}

	if err := h.service.DeleteWithAuth(r.Context(), courseID, userID); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleListCourseEnrollments(w http.ResponseWriter, r *http.Request) {
	courseID := r.PathValue("id")
	if courseID == "" {
		h.writeError(w, errors.New("course_id is required"))
		return
	}

	list, err := h.service.ListByCourseWithAuth(r.Context(), courseID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, list)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, val any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(val)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, ErrNotFound) {
		status = http.StatusNotFound
	} else if err.Error() == "unauthorized" {
		status = http.StatusForbidden
	}
	h.writeJSON(w, status, map[string]string{"error": err.Error()})
}
