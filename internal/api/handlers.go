package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"example.com/pz3-http/internal/storage"
)

type Handlers struct {
	Store *storage.MemoryStore
}

func NewHandlers(store *storage.MemoryStore) *Handlers {
	return &Handlers{Store: store}
}

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		filtered := tasks[:0]
		qq := strings.ToLower(q)
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), qq) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}
	JSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "" && !strings.Contains(ct, "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		BadRequest(w, "title is required")
		return
	}
	if len([]rune(req.Title)) < 3 || len([]rune(req.Title)) > 140 {
		Unprocessable(w, "title length must be between 3 and 140")
		return
	}

	t := h.Store.Create(req.Title)
	JSON(w, http.StatusCreated, t)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.URL.Path, "/tasks/")
	if !ok {
		NotFound(w, "invalid path")
		return
	}
	t, err := h.Store.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			NotFound(w, "task not found")
			return
		}
		JSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unexpected error"})
		return
	}
	JSON(w, http.StatusOK, t)
}

type patchTaskRequest struct {
	Done *bool `json:"done"`
}

func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.URL.Path, "/tasks/")
	if !ok {
		NotFound(w, "invalid path")
		return
	}
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}
	if req.Done == nil {
		Unprocessable(w, "field 'done' is required")
		return
	}
	t, err := h.Store.UpdateDone(id, *req.Done)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			NotFound(w, "task not found")
			return
		}
		JSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unexpected error"})
		return
	}
	JSON(w, http.StatusOK, t)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.URL.Path, "/tasks/")
	if !ok {
		NotFound(w, "invalid path")
		return
	}
	if err := h.Store.Delete(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			NotFound(w, "task not found")
			return
		}
		JSON(w, http.StatusInternalServerError, ErrorResponse{Error: "unexpected error"})
		return
	}
	NoContent(w)
}

func parseID(path string, prefix string) (int64, bool) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	rest := strings.TrimPrefix(path, prefix)
	rest = strings.Trim(rest, "/")
	parts := strings.Split(rest, "/")
	if len(parts) == 0 || parts[0] == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
