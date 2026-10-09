package main

import (
	"errors"
	"io/fs"
	"net/http"
)

type handlers struct{ store *Store }

func newRouter(store *Store, static fs.FS) http.Handler {
	h := &handlers{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/cities", h.cities)
	mux.HandleFunc("POST /api/sessions", h.createSession)
	mux.HandleFunc("POST /api/sessions/{code}/join", h.join)
	mux.HandleFunc("GET /api/sessions/{code}/deck", h.deck)
	mux.HandleFunc("POST /api/sessions/{code}/swipes", h.swipe)
	mux.HandleFunc("GET /api/sessions/{code}/matches", h.matches)
	mux.Handle("GET /", http.FileServerFS(static))
	return mux
}

func (h *handlers) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("OK"))
}

func (h *handlers) cities(w http.ResponseWriter, r *http.Request) {
	cities, err := h.store.Cities(r.Context())
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cities": cities})
}

func (h *handlers) createSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string   `json:"name"`
		City       string   `json:"city"`
		Categories []string `json:"categories"`
	}
	if !decode(w, r, &req) {
		return
	}
	name, ok := cleanName(req.Name)
	if !ok {
		writeError(w, http.StatusBadRequest, "Please enter a name (max 30 characters).")
		return
	}
	cats, ok := cleanCategories(req.Categories)
	if !ok {
		writeError(w, http.StatusBadRequest, "Pick at least one of restaurants, bars or clubs.")
		return
	}
	exists, err := h.store.CityExists(r.Context(), req.City)
	if err != nil {
		serverError(w, err)
		return
	}
	if !exists {
		writeError(w, http.StatusBadRequest, "Unknown city.")
		return
	}
	code, pid, err := h.store.CreateSession(r.Context(), req.City, cats, name)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"code": code, "participantId": pid})
}

func (h *handlers) join(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if !decode(w, r, &req) {
		return
	}
	name, ok := cleanName(req.Name)
	if !ok {
		writeError(w, http.StatusBadRequest, "Please enter a name (max 30 characters).")
		return
	}
	code := normalizeCode(r.PathValue("code"))
	if !validCode(code) {
		writeError(w, http.StatusNotFound, "No session with that code.")
		return
	}
	pid, err := h.store.Join(r.Context(), code, name)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "No session with that code.")
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"code": code, "participantId": pid})
}
