package main

import (
	"errors"
	"net/http"
)

// requireMember resolves the session code and checks the participant is in it.
func (h *handlers) requireMember(w http.ResponseWriter, r *http.Request, pid string) (string, bool) {
	code := normalizeCode(r.PathValue("code"))
	if !validCode(code) || !isUUID(pid) {
		writeError(w, http.StatusNotFound, "Session not found.")
		return "", false
	}
	ok, err := h.store.IsMember(r.Context(), code, pid)
	if err != nil {
		serverError(w, err)
		return "", false
	}
	if !ok {
		writeError(w, http.StatusForbidden, "You're not part of this session. Join it first.")
		return "", false
	}
	return code, true
}

func (h *handlers) deck(w http.ResponseWriter, r *http.Request) {
	pid := r.URL.Query().Get("participant")
	code, ok := h.requireMember(w, r, pid)
	if !ok {
		return
	}
	venues, err := h.store.Deck(r.Context(), code, pid)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"venues": venues})
}

func (h *handlers) swipe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParticipantID string `json:"participantId"`
		VenueID       int    `json:"venueId"`
		Liked         bool   `json:"liked"`
	}
	if !decode(w, r, &req) {
		return
	}
	code, ok := h.requireMember(w, r, req.ParticipantID)
	if !ok {
		return
	}
	match, err := h.store.Swipe(r.Context(), code, req.ParticipantID, req.VenueID, req.Liked)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "Unknown venue.")
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"match": match})
}

func (h *handlers) matches(w http.ResponseWriter, r *http.Request) {
	code := normalizeCode(r.PathValue("code"))
	if !validCode(code) {
		writeError(w, http.StatusNotFound, "Session not found.")
		return
	}
	sum, err := h.store.Summary(r.Context(), code)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "Session not found.")
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sum)
}
