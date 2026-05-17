package config

import (
	"os"

	"go.uber.org/fx"
)

type Config struct {
	Addr        string
	DatabaseURL string
}

func Load() Config {
	cfg := Config{
		Addr:        ":8080",
		DatabaseURL: "./shortener.db.json",
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	return cfg
}

var Module = fx.Module("config", fx.Provide(Load))
