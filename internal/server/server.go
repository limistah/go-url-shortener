package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/limistah/go-url-shortener/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewServer(lc fx.Lifecycle, mux *http.ServeMux, cfg config.Config, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: cfg.Addr, Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("starting server", zap.String("addr", cfg.Addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("server stopped with error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("stopping server")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

var Module = fx.Module(
	"server",
	fx.Provide(zap.NewProduction, NewServer),
	fx.Invoke(func(*http.Server) {}),
)
