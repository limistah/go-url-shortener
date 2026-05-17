package storage

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/limistah/go-url-shortener/internal/config"
	"go.uber.org/fx"
)

type Store interface {
	Save(ctx context.Context, longURL string) (string, error)
	Resolve(ctx context.Context, slug string) (string, bool)
}

type memoryStore struct {
	mu      sync.RWMutex
	prefix  string
	counter int
	bySlug  map[string]string
	byURL   map[string]string
}

func newMemoryStore(prefix string) *memoryStore {
	return &memoryStore{
		prefix: prefix,
		bySlug: make(map[string]string),
		byURL:  make(map[string]string),
	}
}

func (m *memoryStore) Save(_ context.Context, longURL string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if slug, ok := m.byURL[longURL]; ok {
		return slug, nil
	}
	m.counter++
	slug := fmt.Sprintf("%s%d", m.prefix, m.counter)
	m.bySlug[slug] = longURL
	m.byURL[longURL] = slug
	return slug, nil
}

func (m *memoryStore) Resolve(_ context.Context, slug string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.bySlug[slug]
	return v, ok
}

func NewPrimaryStore() Store {
	return newMemoryStore("u")
}

func NewSecondaryStore() Store {
	return newMemoryStore("b")
}

type ActiveStoreIn struct {
	fx.In

	Config    config.Config
	Primary   Store `name:"primary"`
	Secondary Store `name:"secondary"`
}

type mirrorStore struct {
	primary   Store
	secondary Store
}

func (m *mirrorStore) Save(ctx context.Context, longURL string) (string, error) {
	slug, err := m.primary.Save(ctx, longURL)
	if err != nil {
		return "", err
	}
	_, _ = m.secondary.Save(ctx, longURL)
	return slug, nil
}

func (m *mirrorStore) Resolve(ctx context.Context, slug string) (string, bool) {
	if v, ok := m.primary.Resolve(ctx, slug); ok {
		return v, ok
	}
	return m.secondary.Resolve(ctx, slug)
}

func NewActiveStore(in ActiveStoreIn) Store {
	if in.Config.UseSecondary {
		return &mirrorStore{primary: in.Primary, secondary: in.Secondary}
	}
	return in.Primary
}

type validatingStore struct {
	next Store
}

func (v *validatingStore) Save(ctx context.Context, longURL string) (string, error) {
	if strings.TrimSpace(longURL) == "" {
		return "", fmt.Errorf("url is required")
	}
	u, err := url.ParseRequestURI(longURL)
	if err != nil {
		return "", fmt.Errorf("invalid url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("url must start with http or https")
	}
	return v.next.Save(ctx, longURL)
}

func (v *validatingStore) Resolve(ctx context.Context, slug string) (string, bool) {
	return v.next.Resolve(ctx, slug)
}

func DecorateStore(next Store) Store {
	return &validatingStore{next: next}
}

var Module = fx.Module(
	"storage",
	fx.Provide(
		fx.Annotate(NewPrimaryStore, fx.As(new(Store)), fx.ResultTags(`name:"primary"`)),
		fx.Annotate(NewSecondaryStore, fx.As(new(Store)), fx.ResultTags(`name:"secondary"`)),
		NewActiveStore,
	),
)
