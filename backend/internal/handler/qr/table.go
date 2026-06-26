package qr

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type TableManager interface {
	Create(chainID, branchID, label string) (*domain.Table, error)
	ListByBranchID(branchID string) ([]domain.Table, error)
	Update(id, label string) (*domain.Table, error)
	RegenerateQrURL(id string) (*domain.Table, error)
	Delete(id string) error
}

type BranchGetter interface {
	GetByID(id string) (*domain.Branch, error)
}

type Handler struct {
	manager     TableManager
	branchStore BranchGetter
}

func NewHandler(manager TableManager, branchStore BranchGetter) *Handler {
	return &Handler{manager: manager, branchStore: branchStore}
}

type createRequest struct {
	BranchID string `json:"branch_id"`
	Label    string `json:"label"`
}

type updateRequest struct {
	Label string `json:"label"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	branch, err := h.branchStore.GetByID(req.BranchID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "branch not found")
		return
	}

	t, err := h.manager.Create(branch.ChainID, req.BranchID, req.Label)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(w, t)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id required")
		return
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Label == "" {
		response.Error(w, http.StatusBadRequest, "label is required")
		return
	}
	t, err := h.manager.Update(id, req.Label)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, t)
}

func (h *Handler) RegenerateQr(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id required")
		return
	}
	t, err := h.manager.RegenerateQrURL(id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, t)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id required")
		return
	}
	if err := h.manager.Delete(id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, map[string]string{"status": "deleted"})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	branchID := r.URL.Query().Get("branch_id")
	if branchID == "" {
		response.Error(w, http.StatusBadRequest, "branch_id required")
		return
	}
	tables, err := h.manager.ListByBranchID(branchID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, tables)
}
