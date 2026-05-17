package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/limistah/go-url-shortener/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Store interface {
	Save(slug, url string) error
	Lookup(slug string) (string, bool)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

func (s *MemoryStore) Save(slug, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[slug] = url
	return nil
}

func (s *MemoryStore) Lookup(slug string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[slug]
	return url, ok
}

var _ Store = (*MemoryStore)(nil)

type FileStore struct {
	mu   sync.RWMutex
	path string
	data map[string]string
}

func NewFileStore(cfg config.Config) (*FileStore, error) {
	fs := &FileStore{path: cfg.DatabaseURL, data: make(map[string]string)}
	blob, err := os.ReadFile(cfg.DatabaseURL)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fs, nil
		}
		return nil, err
	}
	if len(blob) == 0 {
		return fs, nil
	}
	if err := json.Unmarshal(blob, &fs.data); err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *FileStore) Save(slug, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[slug] = url
	blob, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, blob, 0o644)
}

func (s *FileStore) Lookup(slug string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[slug]
	return url, ok
}

var _ Store = (*FileStore)(nil)

type DualParams struct {
	fx.In

	Hot  Store `name:"hot"`
	Cold Store `name:"cold"`
	Log  *zap.Logger
}

type DualStore struct {
	hot  Store
	cold Store
	log  *zap.Logger
}

func NewDualStore(p DualParams) Store {
	return &DualStore{hot: p.Hot, cold: p.Cold, log: p.Log}
}

func (s *DualStore) Save(slug, url string) error {
	if err := s.hot.Save(slug, url); err != nil {
		return err
	}
	if err := s.cold.Save(slug, url); err != nil {
		s.log.Warn("cold store save failed", zap.Error(err), zap.String("slug", slug))
	}
	return nil
}

func (s *DualStore) Lookup(slug string) (string, bool) {
	if url, ok := s.hot.Lookup(slug); ok {
		return url, true
	}
	return s.cold.Lookup(slug)
}

var _ Store = (*DualStore)(nil)

var Module = fx.Module(
	"storage",
	fx.Provide(
		fx.Annotate(NewMemoryStore, fx.As(new(Store)), fx.ResultTags(`name:"hot"`)),
		fx.Annotate(NewFileStore, fx.As(new(Store)), fx.ResultTags(`name:"cold"`)),
		NewDualStore,
	),
)
