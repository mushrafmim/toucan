package uploads

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"toucan/internal/storage"
)

type Handler struct {
	store storage.Store
}

func NewHandler(store storage.Store) *Handler {
	return &Handler{store: store}
}

type PresignRequest struct {
	Key       string `json:"key"`
	Operation string `json:"operation"`  // "upload" or "download"
	ExpiresIn int    `json:"expires_in"` // seconds
}

type PresignResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func (h *Handler) HandlePresign(w http.ResponseWriter, r *http.Request) {
	var req PresignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	expires := time.Duration(req.ExpiresIn) * time.Second
	if expires == 0 {
		expires = 15 * time.Minute
	}

	var url string
	var err error

	switch req.Operation {
	case "upload":
		url, err = h.store.PresignUpload(r.Context(), req.Key, expires)
	case "download":
		url, err = h.store.PresignDownload(r.Context(), req.Key, expires)
	default:
		http.Error(w, "invalid operation", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to presign: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PresignResponse{
		URL: url,
		Key: req.Key,
	})
}

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Simple multi-part form upload for smaller files or when not using S3 directly
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	key := r.FormValue("key")
	if key == "" {
		key = header.Filename
	}

	if err := h.store.Upload(r.Context(), key, file); err != nil {
		http.Error(w, fmt.Sprintf("failed to upload: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"key": key})
}

func (h *Handler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		key = r.URL.Query().Get("key") // Fallback for some clients
	}

	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	data, err := h.store.Download(r.Context(), key)
	if err != nil {
		if err == storage.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, fmt.Sprintf("failed to download: %v", err), http.StatusInternalServerError)
		return
	}
	defer data.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", key))
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, data); err != nil {
		// Connection might be closed
		return
	}
}
