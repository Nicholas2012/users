package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Nicholas2012/users/internal/api/models"
)

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	// создаём структура куда будем писать
	req := &models.CreateUserRequest{}

	// читаем body в созданную структуру
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.writeErrorBadRequest(w, r, fmt.Errorf("decode request: %w", err))
		return
	}

	if req.Name == "" || req.Surname == "" {
		s.writeErrorBadRequest(w, r, fmt.Errorf("name or surname is empty"))
		return
	}

	// получаем возраст
	slog.Debug("get user age", "name", req.Name)
	age, err := s.ageClient.Age(req.Name)
	if err != nil {
		s.writeErrorInternal(w, r, fmt.Errorf("get user age: %w", err))
		return
	}

	// получаем пол
	slog.Debug("get user gender", "name", req.Name)
	gender, err := s.genderClient.Gender(req.Name)
	if err != nil {
		s.writeErrorInternal(w, r, fmt.Errorf("get gender: %w", err))
		return
	}

	slog.Debug("get user nationality", "name", req.Name)
	nationality, err := s.nationalityClient.Nationality(req.Name)
	if err != nil {
		s.writeErrorInternal(w, r, fmt.Errorf("get nationality: %w", err))
		return
	}

	// делаем вставку в бд
	id, err := s.repo.Insert(req.Name, req.Surname, req.Patronymic, age, gender, nationality)
	if err != nil {
		s.writeErrorInternal(w, r, fmt.Errorf("insert user: %w", err))
		return
	}

	// создаём ответ response
	resp := models.CreateUserResponse{
		ID: id,
	}

	// записываем ответ response
	s.writeJSON(w, r, resp)
}
