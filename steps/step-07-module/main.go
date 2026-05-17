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

type FileStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewFileStore() *FileStore { return &FileStore{data: map[string]string{}} }

func (s *FileStore) Save(slug, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[slug] = longURL
}

func (s *FileStore) Lookup(slug string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[slug]
	return v, ok
}

type DualStoreIn struct {
	fx.In

	Hot  Store `name:"hot"`
	Cold Store `name:"cold"`
}

type DualStore struct {
	hot  Store
	cold Store
}

func NewDualStore(in DualStoreIn) Store { return &DualStore{hot: in.Hot, cold: in.Cold} }

func (d *DualStore) Save(slug, longURL string) {
	d.hot.Save(slug, longURL)
	d.cold.Save(slug, longURL)
}

func (d *DualStore) Lookup(slug string) (string, bool) {
	if v, ok := d.hot.Lookup(slug); ok {
		return v, true
	}
	return d.cold.Lookup(slug)
}

type Route interface {
	http.Handler
	Pattern() string
}

type ShortenParams struct {
	fx.In

	Store  Store
	Config Config
}

type ShortenHandler struct {
	store   Store
	cfg     Config
	counter atomic.Uint64
}

func NewShortenHandler(p ShortenParams) *ShortenHandler {
	return &ShortenHandler{store: p.Store, cfg: p.Config}
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
}

type RedirectParams struct {
	fx.In

	Store Store
}

type RedirectHandler struct{ store Store }

func NewRedirectHandler(p RedirectParams) *RedirectHandler { return &RedirectHandler{store: p.Store} }

func (h *RedirectHandler) Pattern() string { return "GET /{slug}" }

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	v, ok := h.store.Lookup(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
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

var ConfigModule = fx.Module("config", fx.Provide(Load))

var StorageModule = fx.Module("storage", fx.Provide(
	fx.Annotate(NewMemoryStore, fx.As(new(Store)), fx.ResultTags(`name:"hot"`)),
	fx.Annotate(NewFileStore, fx.As(new(Store)), fx.ResultTags(`name:"cold"`)),
	NewDualStore,
))

var APIModule = fx.Module("api", fx.Provide(
	fx.Annotate(NewShortenHandler, fx.As(new(Route)), fx.ResultTags(`group:"routes"`)),
	fx.Annotate(NewRedirectHandler, fx.As(new(Route)), fx.ResultTags(`group:"routes"`)),
	fx.Annotate(NewMux, fx.ParamTags(`group:"routes"`)),
))

var ServerModule = fx.Module("server", fx.Provide(zap.NewProduction, NewServer), fx.Invoke(func(*http.Server) {}))

func main() {
	fx.New(ConfigModule, StorageModule, APIModule, ServerModule).Run()
}
