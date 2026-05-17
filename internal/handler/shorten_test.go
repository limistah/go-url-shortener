package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/limistah/go-url-shortener/internal/config"
	"github.com/limistah/go-url-shortener/internal/handler"
	"github.com/limistah/go-url-shortener/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func testConfig() config.Config {
	return config.Config{Addr: ":8080", DatabaseURL: ""}
}

func TestPopulateShortenHandler(t *testing.T) {
	var h *handler.ShortenHandler

	app := fxtest.New(t,
		fx.Provide(
			testConfig,
			zap.NewNop,
			fx.Annotate(storage.NewMemoryStore, fx.As(new(storage.Store))),
			handler.NewShortenHandler,
		),
		fx.Populate(&h),
	)
	app.RequireStart()
	defer app.RequireStop()

	if h == nil {
		t.Fatal("expected populated shorten handler")
	}
}

type fakeStore struct{}

func (fakeStore) Save(string, string) error    { return nil }
func (fakeStore) Lookup(string) (string, bool) { return "https://example.com", true }

func TestReplaceStore(t *testing.T) {
	var h *handler.ShortenHandler

	app := fxtest.New(t,
		fx.Provide(
			testConfig,
			zap.NewNop,
			fx.Annotate(storage.NewMemoryStore, fx.As(new(storage.Store))),
			handler.NewShortenHandler,
		),
		fx.Replace(storage.Store(fakeStore{})),
		fx.Populate(&h),
	)
	app.RequireStart()
	defer app.RequireStop()

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(`{"url":"https://example.com"}`))
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}
}
