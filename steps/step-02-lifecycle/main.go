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

type ShortenHandler struct {
	store   *MemoryStore
	logger  *zap.Logger
	cfg     Config
	counter atomic.Uint64
}

func NewShortenHandler(store *MemoryStore, logger *zap.Logger, cfg Config) *ShortenHandler {
	return &ShortenHandler{store: store, logger: logger, cfg: cfg}
}

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
}

type RedirectHandler struct {
	store  *MemoryStore
	logger *zap.Logger
}

func NewRedirectHandler(store *MemoryStore, logger *zap.Logger) *RedirectHandler {
	return &RedirectHandler{store: store, logger: logger}
}

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

func NewMux(s *ShortenHandler, r *RedirectHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("POST /shorten", s)
	mux.Handle("GET /{slug}", r)
	return mux
}

func NewServer(lc fx.Lifecycle, mux *http.ServeMux, cfg Config, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: cfg.Addr, Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("starting server", zap.String("addr", cfg.Addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("listen failed", zap.Error(err))
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

func main() {
	fx.New(
		fx.Provide(
			Load,
			zap.NewProduction,
			NewMemoryStore,
			NewShortenHandler,
			NewRedirectHandler,
			NewMux,
			NewServer,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
