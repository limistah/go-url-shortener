package app_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/limistah/go-url-shortener/internal/api/handlers"
	"github.com/limistah/go-url-shortener/internal/app"
	"github.com/limistah/go-url-shortener/internal/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestPopulateHandlerWithFxTest(t *testing.T) {
	var h *handlers.URLHandler

	fxApp := fxtest.New(t,
		app.Module,
		fx.NopLogger,
		fx.Populate(&h),
	)
	fxApp.RequireStart()
	defer fxApp.RequireStop()

	if h == nil {
		t.Fatal("expected populated handler")
	}
}

func TestReplaceStoreForHandler(t *testing.T) {
	var h *handlers.URLHandler

	fxApp := fxtest.New(t,
		app.Module,
		fx.NopLogger,
		fx.Replace(config.Config{Addr: ":18080", BaseURL: "http://short.local", UseSecondary: false}),
		fx.Populate(&h),
	)
	fxApp.RequireStart()
	defer fxApp.RequireStop()

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(`{"url":"https://example.org"}`))
	res := httptest.NewRecorder()

	h.HandleShorten(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, res.Code)
	}
	if got := res.Body.String(); !bytes.Contains([]byte(got), []byte("http://short.local/")) {
		t.Fatalf("expected replaced base URL in response, got: %s", got)
	}
}
