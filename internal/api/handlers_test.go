package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"example.com/pz3-http/internal/storage"
)

// --- Helpers ---

func mustUnmarshal[T any](t *testing.T, b []byte) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("invalid json: %v\nbody: %s", err, string(b))
	}
	return v
}

// --- Tests ---

// Создание задачи: ожидаем 201 и корректный JSON
func TestCreateTask(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	body := bytes.NewBufferString(`{"title":"Buy milk"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body=%s", rec.Code, rec.Body.String())
	}
	task := mustUnmarshal[storage.Task](t, rec.Body.Bytes())
	if task.Title != "Buy milk" || task.Done {
		t.Fatalf("unexpected task: %+v", task)
	}
}

// Список задач: после Create должен вернуться 1 элемент
func TestListTasks(t *testing.T) {
	store := storage.NewMemoryStore()
	store.Create("Buy milk")
	h := NewHandlers(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	h.ListTasks(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	tasks := mustUnmarshal[[]storage.Task](t, rec.Body.Bytes())
	if len(tasks) != 1 || tasks[0].Title != "Buy milk" {
		t.Fatalf("unexpected tasks: %+v", tasks)
	}
}

// Получение по id: OK
func TestGetTaskOK(t *testing.T) {
	store := storage.NewMemoryStore()
	created := store.Create("Buy milk")
	h := NewHandlers(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+strconv.FormatInt(created.ID, 10), nil)
	rec := httptest.NewRecorder()

	h.GetTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", rec.Code, rec.Body.String())
	}
	got := mustUnmarshal[storage.Task](t, rec.Body.Bytes())
	if got.ID != created.ID || got.Title != "Buy milk" {
		t.Fatalf("unexpected task: %+v", got)
	}
}

// Валидация: короткий title -> 422
func TestCreateValidationLength(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	body := bytes.NewBufferString(`{"title":"ok"}`) // < 3 символов
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d; body=%s", rec.Code, rec.Body.String())
	}
}

// Неверный Content-Type -> 400 (в хендлере это проверяется)
func TestCreateWrongContentType(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	body := bytes.NewBufferString(`{"title":"Buy milk"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "text/plain") // специально неверный
	rec := httptest.NewRecorder()

	h.CreateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d; body=%s", rec.Code, rec.Body.String())
	}
}

// PATCH done=true: ожидаем 200 и Done=true
func TestPatchTaskDone(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)
	created := store.Create("Buy milk")

	body := bytes.NewBufferString(`{"done":true}`)
	req := httptest.NewRequest(http.MethodPatch, "/tasks/"+strconv.FormatInt(created.ID, 10), body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.PatchTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", rec.Code, rec.Body.String())
	}
	got := mustUnmarshal[storage.Task](t, rec.Body.Bytes())
	if !got.Done {
		t.Fatalf("expected done=true, got %+v", got)
	}
}

// DELETE: 204, затем GET -> 404
func TestDeleteTask(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)
	created := store.Create("Buy milk")

	// delete
	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+strconv.FormatInt(created.ID, 10), nil)
	rec := httptest.NewRecorder()
	h.DeleteTask(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d; body=%s", rec.Code, rec.Body.String())
	}

	// get after delete -> 404
	req2 := httptest.NewRequest(http.MethodGet, "/tasks/"+strconv.FormatInt(created.ID, 10), nil)
	rec2 := httptest.NewRecorder()
	h.GetTask(rec2, req2)

	if rec2.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d; body=%s", rec2.Code, rec2.Body.String())
	}
}

// Ошибка пути: /tasks/abc -> 404 (invalid path)
func TestGetTaskInvalidPath(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	rec := httptest.NewRecorder()
	h.GetTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for invalid path, got %d; body=%s", rec.Code, rec.Body.String())
	}
}

// Не найдено: корректный id, которого нет -> 404
func TestGetTaskNotFound(t *testing.T) {
	store := storage.NewMemoryStore()
	h := NewHandlers(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks/9999", nil)
	rec := httptest.NewRecorder()
	h.GetTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d; body=%s", rec.Code, rec.Body.String())
	}
}
