package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct{ Addr string }

func Load() Config { return Config{Addr: ":8080"} }

type Store interface {
	Save(slug, longURL string)
	Lookup(slug string) (string, bool)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{data: map[string]string{}} }

func (s *MemoryStore) Save(slug, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[slug] = longURL
}

func (s *MemoryStore) Lookup(slug string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[slug]
	return v, ok
}

var _ Store = (*MemoryStore)(nil)

type Route interface {
	http.Handler
	Pattern() string
}

type ShortenParams struct {
	fx.In

	Store  Store
	Logger *zap.Logger
	Config Config
}

type ShortenHandler struct {
	store   Store
	logger  *zap.Logger
	cfg     Config
	counter atomic.Uint64
}

func NewShortenHandler(p ShortenParams) *ShortenHandler {
	return &ShortenHandler{store: p.Store, logger: p.Logger, cfg: p.Config}
}

func (h *ShortenHandler) Pattern() string { return "POST /shorten" }

func (h *ShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	u, err := url.ParseRequestURI(body.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	slug := fmt.Sprintf("u%x", h.counter.Add(1))
	h.store.Save(slug, body.URL)
	base := "http://localhost" + strings.TrimSpace(h.cfg.Addr)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"slug": slug, "short_url": base + "/" + slug})
	h.logger.Info("shortened", zap.String("slug", slug))
}

type RedirectParams struct {
	fx.In

	Store  Store
	Logger *zap.Logger
}

type RedirectHandler struct {
	store  Store
	logger *zap.Logger
}

func NewRedirectHandler(p RedirectParams) *RedirectHandler {
	return &RedirectHandler{store: p.Store, logger: p.Logger}
}

func (h *RedirectHandler) Pattern() string { return "GET /{slug}" }

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	v, ok := h.store.Lookup(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	h.logger.Info("redirect", zap.String("slug", slug), zap.String("url", v))
	http.Redirect(w, r, v, http.StatusFound)
}

func NewMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

func NewServer(lc fx.Lifecycle, mux *http.ServeMux, cfg Config, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: cfg.Addr, Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("listen failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error { return srv.Shutdown(ctx) },
	})
	return srv
}

var Module = fx.Options(
	fx.Provide(
		Load,
		zap.NewProduction,
		fx.Annotate(NewMemoryStore, fx.As(new(Store))),
		fx.Annotate(NewShortenHandler, fx.As(new(Route)), fx.ResultTags(`group:"routes"`)),
		fx.Annotate(NewRedirectHandler, fx.As(new(Route)), fx.ResultTags(`group:"routes"`)),
		fx.Annotate(NewMux, fx.ParamTags(`group:"routes"`)),
		NewServer,
	),
	fx.Invoke(func(*http.Server) {}),
)

func main() {
	fx.New(Module).Run()
}
