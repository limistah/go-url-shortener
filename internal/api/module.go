package api

import (
	"net/http"

	"github.com/limistah/go-url-shortener/internal/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewMux(routes []handler.Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range routes {
		mux.Handle(r.Pattern(), r)
	}
	return mux
}

var Module = fx.Module(
	"api",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.With(zap.String("module", "api"))
	}),
	fx.Provide(
		fx.Annotate(handler.NewShortenHandler, fx.As(new(handler.Route)), fx.ResultTags(`group:"routes"`)),
		fx.Annotate(handler.NewRedirectHandler, fx.As(new(handler.Route)), fx.ResultTags(`group:"routes"`)),
		fx.Annotate(NewMux, fx.ParamTags(`group:"routes"`)),
	),
)
