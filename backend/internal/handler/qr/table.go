package qr

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type TableManager interface {
	Create(branchID, label string) (*domain.Table, error)
	ListByBranchID(branchID string) ([]domain.Table, error)
	Delete(id string) error
}

type Handler struct {
	manager TableManager
}

func NewHandler(manager TableManager) *Handler {
	return &Handler{manager: manager}
}

type createRequest struct {
	BranchID string `json:"branch_id"`
	Label    string `json:"label"`
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
	t, err := h.manager.Create(req.BranchID, req.Label)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, t)
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
	response.JSON(w, http.StatusOK, map[string]string{"status": "deleted"})
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
	response.JSON(w, http.StatusOK, tables)
}
