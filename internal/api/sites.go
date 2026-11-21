package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kaherro/MultiPing/internal/domain"
)

type SiteStore interface {
	Create(ctx context.Context, s *domain.Site) error
	List(ctx context.Context) ([]domain.Site, error)
	Delete(ctx context.Context, id int64) error
}

type SitesHandler struct {
	store SiteStore
	logger *slog.Logger
}

func NewSitesHandler(store SiteStore, logger *slog.Logger) *SitesHandler {
	return &SitesHandler{store: store, logger: logger}
}

func (h *SitesHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST sites", h.create); 
	mux.HandleFunc("GET sites", h.list); 
	mux.HandleFunc("DELETE sites", h.delete);
}

func (h *SitesHandler) create(w http.ResponseWriter, r *http.Request) {
	var s domain.Site
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return 
	}
	
	u, err := url.ParseRequestURI(s.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		writeError(w, http.StatusBadRequest, "url must be a valid http(s) url")
		return
	}
	if s.IntervalSec == 0 {
		s.IntervalSec = 60
	}
	if s.TimeoutMS == 0 {
		s.TimeoutMS = 5000
	}
	if s.IntervalSec < 10 || s.TimeoutMS < 100 {
		writeError(w, http.StatusBadRequest, "interval_sec must be >= 10 and timeout_ms >= 100")
		return
	}
	if err := h.store.Create(r.Context(), &s); err != nil {
		h.logger.Error("create site", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

func (h *SitesHandler) list(w http.ResponseWriter, r *http.Request) {
	sites, err := h.store.List(r.Context())
	if err != nil {
		h.logger.Error("list sites", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, sites)
}

func (h *SitesHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.store.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "site not found")
			return
		}
		h.logger.Error("delete site", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}