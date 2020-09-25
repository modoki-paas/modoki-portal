package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func (h *handler) saveSession(s *sessions.Session, rw http.ResponseWriter, req *http.Request) bool {
	if err := s.Save(req, rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return false
	}

	return true
}

func (h *handler) loadSession(rw http.ResponseWriter, req *http.Request) *sessions.Session {
	session, err := h.store.Get(req, cookieKey)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)

		return nil
	}

	return session
}
