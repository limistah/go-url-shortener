package app

import (
	"github.com/limistah/go-url-shortener/internal/api/handlers"
	"github.com/limistah/go-url-shortener/internal/api/routes"
	"github.com/limistah/go-url-shortener/internal/config"
	"github.com/limistah/go-url-shortener/internal/server"
	"github.com/limistah/go-url-shortener/internal/storage"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"url-shortener",
	config.Module,
	storage.Module,
	handlers.Module,
	routes.Module,
	server.Module,
	fx.Decorate(storage.DecorateStore),
)

func New() *fx.App {
	return fx.New(Module)
}
