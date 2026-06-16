package chain

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type CreateBranchRequest struct {
	ChainID     string  `json:"chain_id"`
	Name        string  `json:"name"`
	Address     string  `json:"address,omitempty"`
	ContactInfo string  `json:"contact_info,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

type BranchCreator interface {
	Create(chainID, name, address, contactInfo string, parentID *string) (*domain.Branch, error)
	ListByChainID(chainID string) ([]domain.Branch, error)
}

type BranchHandler struct {
	creator BranchCreator
}

func NewBranchHandler(creator BranchCreator) *BranchHandler {
	return &BranchHandler{creator: creator}
}

func (h *BranchHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	branch, err := h.creator.Create(req.ChainID, req.Name, req.Address, req.ContactInfo, req.ParentID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, branch)
}

func (h *BranchHandler) ListByChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id is required")
		return
	}

	branches, err := h.creator.ListByChainID(chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, branches)
}
