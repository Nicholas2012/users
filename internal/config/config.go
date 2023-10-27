package config

import (
	"os"
	"strings"
)

type Mode string

const (
	Dev  Mode = "dev"  // режим разработки
	Prod Mode = "prod" // режим продакшена
)

type Config struct {
	Listen string // адрес, на котором будет запущен сервер
	DB     string // строка подключения к базе данных
	Mode   Mode   // режим работы приложения
}

// IsDev возвращает true, если приложение запущено в режиме разработки
func (c Config) IsDev() bool {
	return c.Mode == Dev
}

func Load() Config {
	return Config{
		Listen: getString("LISTEN", ":8080"),
		DB:     getString("DB", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		Mode:   getMode("MODE", Dev),
	}
}

func getString(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getMode(key string, def Mode) Mode {
	switch v := Mode(strings.ToLower(getString(key, string(def)))); v {
	case Dev:
		return Dev
	case Prod:
		return Prod
	default:
		return def
	}
}
