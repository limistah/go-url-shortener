package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/limistah/go-url-shortener/internal/config"
	"go.uber.org/fx"
)

func NewMux() *mux.Router {
	return mux.NewRouter()
}

func NewHTTPServer(cfg config.Config, mux *mux.Router) *http.Server {
	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

type LifecycleParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Server    *http.Server
}

func RegisterLifecycle(p LifecycleParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := p.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("http server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return p.Server.Shutdown(ctx)
		},
	})
}

var Module = fx.Module(
	"server",
	fx.Provide(NewMux, NewHTTPServer),
	fx.Invoke(RegisterLifecycle),
)
