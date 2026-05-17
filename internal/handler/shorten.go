package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/limistah/go-url-shortener/internal/config"
	"github.com/limistah/go-url-shortener/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ShortenParams struct {
	fx.In

	Store  storage.Store
	Logger *zap.Logger
	Config config.Config
}

type ShortenHandler struct {
	store   storage.Store
	logger  *zap.Logger
	cfg     config.Config
	counter atomic.Uint64
}

func NewShortenHandler(p ShortenParams) *ShortenHandler {
	return &ShortenHandler{store: p.Store, logger: p.Logger, cfg: p.Config}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Slug     string `json:"slug"`
	ShortURL string `json:"short_url"`
}

func (h *ShortenHandler) Pattern() string { return "POST /shorten" }

func (h *ShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.URL) == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}
	u, err := url.ParseRequestURI(req.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}
	n := h.counter.Add(1)
	slug := fmt.Sprintf("u%x", n)
	if err := h.store.Save(slug, req.URL); err != nil {
		h.logger.Error("failed to save shortened url", zap.Error(err), zap.String("slug", slug))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	baseURL := "http://localhost" + h.cfg.Addr
	res := shortenResponse{Slug: slug, ShortURL: baseURL + "/" + slug}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
	}
}
