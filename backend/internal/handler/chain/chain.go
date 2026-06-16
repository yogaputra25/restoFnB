package chain

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type RegisterRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ContactInfo string `json:"contact_info,omitempty"`
}

type Registrar interface {
	Register(name, slug, contactInfo string) (*domain.Chain, error)
}

type Handler struct {
	registrar Registrar
}

func NewHandler(registrar Registrar) *Handler {
	return &Handler{registrar: registrar}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	chain, err := h.registrar.Register(req.Name, req.Slug, req.ContactInfo)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, chain)
}
