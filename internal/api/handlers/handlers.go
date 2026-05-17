package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/limistah/go-url-shortener/internal/config"
	"github.com/limistah/go-url-shortener/internal/storage"
	"go.uber.org/fx"
)

type URLHandler struct {
	store   storage.Store
	baseURL string
}

type Params struct {
	fx.In

	Store  storage.Store
	Config config.Config
}

func NewURLHandler(p Params) *URLHandler {
	return &URLHandler{store: p.Store, baseURL: strings.TrimRight(p.Config.BaseURL, "/")}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Slug     string `json:"slug"`
	ShortURL string `json:"short_url"`
}

func (h *URLHandler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	slug, err := h.store.Save(r.Context(), req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(shortenResponse{Slug: slug, ShortURL: h.baseURL + "/" + slug}); err != nil {
		log.Printf("encode shorten response failed: %v", err)
	}
}

func (h *URLHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	longURL, ok := h.store.Resolve(r.Context(), slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

var Module = fx.Module("handlers", fx.Provide(NewURLHandler))
