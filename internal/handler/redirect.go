package handler

import (
	"net/http"

	"github.com/limistah/go-url-shortener/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RedirectParams struct {
	fx.In

	Store  storage.Store
	Logger *zap.Logger
}

type RedirectHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewRedirectHandler(p RedirectParams) *RedirectHandler {
	return &RedirectHandler{store: p.Store, logger: p.Logger}
}

func (h *RedirectHandler) Pattern() string { return "GET /{slug}" }

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	longURL, ok := h.store.Lookup(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	h.logger.Info("redirect", zap.String("slug", slug), zap.String("url", longURL))
	http.Redirect(w, r, longURL, http.StatusFound)
}
