package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Nicholas2012/users/internal/api/models"
	"github.com/Nicholas2012/users/internal/repository"
	"github.com/go-chi/chi/v5"
)

type Ager interface {
	Age(name string) (int, error)
}

type Genders interface {
	Gender(name string) (string, error)
}

type Nationalities interface {
	Nationality(name string) (string, error)
}

type UsersRepository interface {
	Get(id int) (*repository.User, error)
	List(opts repository.ListOpts) (*repository.UserList, error)
	Insert(name string, surname string, patronymic string, age int, gender string, nationality string) (int, error)
	Update(id int, name string, surname string, patronymic string, age int, gender string, nationality string) (int, error)
	Delete(id int) error
}

type Service struct {
	repo              UsersRepository
	ageClient         Ager
	genderClient      Genders
	nationalityClient Nationalities
}

func NewService(r UsersRepository, a Ager, g Genders, n Nationalities) *Service {
	return &Service{
		repo:              r,
		ageClient:         a,
		genderClient:      g,
		nationalityClient: n,
	}
}

func (s *Service) RegisterRoutes(r *chi.Mux) {
	r.Get("/users", s.ListUsers)
	r.Post("/users", s.CreateUser)
	r.Get("/users/{id}", s.GetUser)
	r.Delete("/users/{id}", s.DeleteUser)
	r.Put("/users/{id}", s.Update)
}

func (s *Service) writeJSON(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.writeErrorInternal(w, r, err)
	}
}

func (s *Service) writeError(w http.ResponseWriter, r *http.Request, status int, err error) {
	slog.Error("request error", "err", err,
		"url", r.URL.String(),
		"method", r.Method,
	)

	w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()}); err != nil {
		slog.Error("write error response", "err", err)
	}
}

func (s *Service) writeErrorBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	s.writeError(w, r, http.StatusBadRequest, err)
}

func (s *Service) writeErrorInternal(w http.ResponseWriter, r *http.Request, err error) {
	s.writeError(w, r, http.StatusInternalServerError, err)
}
