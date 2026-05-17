package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/limistah/go-url-shortener/internal/api/handlers"
	"go.uber.org/fx"
)

type Route interface {
	Register(r *mux.Router)
}

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func (r route) Register(m *mux.Router) {
	m.HandleFunc(r.path, r.handler).Methods(r.method)
}

type Out struct {
	fx.Out

	Shorten Route `group:"routes"`
	Lookup  Route `group:"routes"`
}

func NewRoutes(h *handlers.URLHandler) Out {
	return Out{
		Shorten: route{method: http.MethodPost, path: "/shorten", handler: h.HandleShorten},
		Lookup:  route{method: http.MethodGet, path: "/{slug}", handler: h.HandleRedirect},
	}
}

type RegisterParams struct {
	fx.In

	Mux    *mux.Router
	Routes []Route `group:"routes"`
}

func RegisterAll(p RegisterParams) {
	for _, rt := range p.Routes {
		rt.Register(p.Mux)
	}
}

var Module = fx.Module(
	"routes",
	fx.Provide(NewRoutes),
	fx.Invoke(RegisterAll),
)
