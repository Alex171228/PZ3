package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/pz3-http/internal/storage"
)

func setup() (*Handlers, http.Handler) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	return h, mux
}

func TestCreateTaskAndList(t *testing.T) {
	_, mux := setup()

	body := bytes.NewBufferString(`{"title":"Купить молоко"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body=%s", rec.Code, rec.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if created["title"] != "Купить молоко" {
		t.Fatalf("unexpected title: %v", created["title"])
	}

	req2 := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec2.Code)
	}
	var arr []map[string]any
	if err := json.Unmarshal(rec2.Body.Bytes(), &arr); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 task, got %d", len(arr))
	}
}

func TestCreateValidationLength(t *testing.T) {
	_, mux := setup()

	body := bytes.NewBufferString(`{"title":"ok"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d; body=%s", rec.Code, rec.Body.String())
	}
}
