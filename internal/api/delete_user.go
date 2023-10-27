package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		s.writeErrorBadRequest(w, r, err)
		return
	}

	if err := s.repo.Delete(idInt); err != nil {
		s.writeErrorInternal(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
