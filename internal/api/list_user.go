package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Nicholas2012/users/internal/api/models"
)

func (s *Service) ListUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		pageInt = 1
	}

	userList, err := s.repo.List(pageInt)
	if err != nil {
		s.writeErrorInternal(w, r, err)
		return
	}

	slog.Debug("list users",
		"total", userList.Count,
		"page", userList.Page,
		"pages", userList.Pages,
		"users", len(userList.Users),
	)

	response := models.ListUsersResponse{
		Page:  userList.Page,
		Pages: userList.Pages,
		Users: make([]models.User, len(userList.Users)),
	}
	for i, u := range userList.Users {
		response.Users[i] = models.User(u)
	}

	s.writeJSON(w, r, response)
}
