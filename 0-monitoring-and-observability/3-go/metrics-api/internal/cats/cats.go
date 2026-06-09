// Package cats implements the demo /cats REST resource (list + create) so the
// service emits real routes that produce metrics. State is in-memory: the store
// is seeded with one row (Tom) so GET /cats is meaningful from the first call.
package cats

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Cat is the demo resource returned by GET /cats and created by POST /cats.
type Cat struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type store struct {
	mu     sync.Mutex
	cats   []Cat
	nextID int
}

var s = &store{
	cats:   []Cat{{ID: 1, Name: "Tom", Age: 3}},
	nextID: 2,
}

// createRequest mirrors CreateCatDto: name is required, age must be an integer.
type createRequest struct {
	Name *string `json:"name"`
	Age  *int    `json:"age"`
}

// Handler routes GET (list) and POST (create) on /cats, matching the shared
// API contract: GET 200 returns the list, POST with a missing name returns 400.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleList(w)
		case http.MethodPost:
			handleCreate(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, []string{"method not allowed"})
		}
	})
}

func handleList(w http.ResponseWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, s.cats)
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var body createRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, []string{"invalid JSON body"})
		return
	}
	// Validation: name must be present and non-empty (mirrors @IsNotEmpty).
	if body.Name == nil || *body.Name == "" {
		writeError(w, http.StatusBadRequest, []string{"name should not be empty"})
		return
	}
	if body.Age == nil {
		writeError(w, http.StatusBadRequest, []string{"age must be an integer number"})
		return
	}
	s.mu.Lock()
	cat := Cat{ID: s.nextID, Name: *body.Name, Age: *body.Age}
	s.nextID++
	s.cats = append(s.cats, cat)
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, cat)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, messages []string) {
	writeJSON(w, status, map[string]any{
		"statusCode": status,
		"message":    messages,
		"error":      http.StatusText(status),
	})
}
