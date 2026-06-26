package theme

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type UpsertRequest struct {
	ChainID      string `json:"chain_id"`
	PrimaryColor string `json:"primary_color"`
	AccentColor  string `json:"accent_color"`
	LogoURL      string `json:"logo_url,omitempty"`
}

type ThemeManager interface {
	Upsert(chainID, primaryColor, accentColor, logoURL string) (*domain.Theme, error)
	GetByChainID(chainID string) (*domain.Theme, error)
}

type Handler struct {
	manager ThemeManager
}

func NewHandler(manager ThemeManager) *Handler {
	return &Handler{manager: manager}
}

func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req UpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	theme, err := h.manager.Upsert(req.ChainID, req.PrimaryColor, req.AccentColor, req.LogoURL)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, theme)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id is required")
		return
	}

	theme, err := h.manager.GetByChainID(chainID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "theme not found")
		return
	}

	response.Success(w, theme)
}
