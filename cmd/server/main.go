package main

import (
	"go.uber.org/fx"

	"github.com/limistah/go-url-shortener/internal/api"
	"github.com/limistah/go-url-shortener/internal/config"
	"github.com/limistah/go-url-shortener/internal/server"
	"github.com/limistah/go-url-shortener/internal/storage"
)

func main() {
	fx.New(
		config.Module,
		storage.Module,
		api.Module,
		server.Module,
	).Run()
}
