package config

import (
	"os"

	"go.uber.org/fx"
)

type Config struct {
	Addr         string
	BaseURL      string
	UseSecondary bool
}

func New() Config {
	cfg := Config{
		Addr:         ":8080",
		BaseURL:      "http://localhost:8080",
		UseSecondary: false,
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if os.Getenv("USE_SECONDARY_STORE") == "1" {
		cfg.UseSecondary = true
	}
	return cfg
}

var Module = fx.Module("config", fx.Provide(New))
