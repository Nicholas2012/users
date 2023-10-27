package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		s.writeErrorBadRequest(w, r, err)
		return
	}

	user, err := s.repo.Get(idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("user not found", "id", idInt)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		s.writeErrorInternal(w, r, err)
		return
	}

	s.writeJSON(w, r, user)
}
