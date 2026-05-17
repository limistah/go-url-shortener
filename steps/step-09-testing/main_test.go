package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestPopulateShortenHandler(t *testing.T) {
	var h *ShortenHandler

	app := fxtest.New(t,
		fx.Provide(
			Load,
			zap.NewNop,
			fx.Annotate(NewMemoryStore, fx.As(new(Store))),
			NewShortenHandler,
		),
		fx.Populate(&h),
	)
	app.RequireStart()
	defer app.RequireStop()

	if h == nil {
		t.Fatal("expected shorten handler")
	}
}

type fakeStore struct{}

func (fakeStore) Save(string, string)          {}
func (fakeStore) Lookup(string) (string, bool) { return "https://example.com", true }

func TestReplaceStore(t *testing.T) {
	var h *ShortenHandler

	app := fxtest.New(t,
		fx.Provide(
			Load,
			zap.NewNop,
			fx.Annotate(NewMemoryStore, fx.As(new(Store))),
			NewShortenHandler,
		),
		fx.Replace(Store(fakeStore{})),
		fx.Populate(&h),
	)
	app.RequireStart()
	defer app.RequireStop()

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(`{"url":"https://example.com"}`))
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
}
