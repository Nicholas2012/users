package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Nicholas2012/users/internal/api/models"
	"github.com/Nicholas2012/users/internal/repository"
)

func (s *Service) ListUsers(w http.ResponseWriter, r *http.Request) {
	var (
		pageStr  = r.URL.Query().Get("page")
		limitStr = r.URL.Query().Get("limit")
		ageStr   = r.URL.Query().Get("age")
		nameStr  = r.URL.Query().Get("name")
	)

	opts := repository.ListOpts{
		Page:  1,
		Limit: 20,
		Name:  nameStr,
	}

	if v, err := strconv.Atoi(pageStr); err == nil {
		opts.Page = v
	}

	if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
		opts.Limit = v
	}

	if v, err := strconv.Atoi(ageStr); err == nil {
		opts.Age = v
	}

	userList, err := s.repo.List(opts)
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
