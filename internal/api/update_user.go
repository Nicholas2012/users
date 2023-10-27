package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Nicholas2012/users/internal/api/models"
	"github.com/go-chi/chi/v5"
)

func (s *Service) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		s.writeErrorBadRequest(w, r, err)
		return
	}

	// создаём структура куда будем писать
	req := &models.UpdateUserRequest{}

	// читаем body в созданную структуру
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.writeErrorBadRequest(w, r, fmt.Errorf("decode request: %w", err))
		return
	}

	if req.Name == "" || req.Surname == "" {
		s.writeErrorBadRequest(w, r, fmt.Errorf("name or surname is empty"))
		return
	}

	// делаем вставку в бд
	id, err := s.repo.Update(idInt, req.Name, req.Surname, req.Patronymic, req.Age, req.Gender, req.Nationality)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("user not found", "id", idInt)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		s.writeErrorInternal(w, r, fmt.Errorf("insert user: %w", err))
		return
	}

	// создаём ответ response
	resp := models.UpdateUserResponse{
		ID: id,
	}

	// записываем ответ response
	s.writeJSON(w, r, resp)
}
