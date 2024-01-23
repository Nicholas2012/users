package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Nicholas2012/users/internal/api"
	"github.com/Nicholas2012/users/internal/config"
	"github.com/Nicholas2012/users/internal/pkg/agify"
	"github.com/Nicholas2012/users/internal/pkg/genderize"
	m "github.com/Nicholas2012/users/internal/pkg/middleware"
	"github.com/Nicholas2012/users/internal/pkg/nationalize"
	"github.com/Nicholas2012/users/internal/repository"
	"github.com/Nicholas2012/users/migrations"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	setupLogger(cfg)

	slog.Debug("start app", slog.String("mode", string(cfg.Mode)))

	db, err := setupDB(cfg)
	if err != nil {
		slog.Error("setup db connection", "err", err)
		os.Exit(1)
	}

	storage := repository.New(db) // создали клиента к базе данных

	ageClient := agify.New(cfg.AgeAPI)
	genderClient := genderize.New(cfg.GenderAPI)
	nationalityClient := nationalize.New(cfg.NationalityAPI)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(m.Logger)

	service := api.NewService(storage, ageClient, genderClient, nationalityClient)
	// на основе клиента создаем сервис
	service.RegisterRoutes(r)

	slog.Info("start listen", "addr", cfg.Listen)
	if err := http.ListenAndServe(cfg.Listen, r); err != nil {
		slog.Error("listen and serve", "err", err)
		os.Exit(1)
	}
}

func setupDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	goose.SetBaseFS(migrations.Migrations)
	if err := goose.Up(db, "."); err != nil {
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	return db, nil
}

func setupLogger(cfg config.Config) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if !cfg.IsDev() {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(logger)
}
