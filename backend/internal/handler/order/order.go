package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type OrderService interface {
	Create(chainID, branchID, orderType string, tableID, customerID, guestID *string, customerName string, items []domain.OrderItem) (*domain.Order, error)
	UpdateStatus(id, status string) (*domain.Order, error)
	Pay(id, method string) (*domain.Order, error)
	ListByBranchID(branchID string, page, limit int) ([]domain.Order, int, error)
	ListByChainID(chainID string, page, limit int) ([]domain.Order, int, error)
}

type Handler struct {
	svc OrderService
}

func NewHandler(svc OrderService) *Handler {
	return &Handler{svc: svc}
}

type createRequest struct {
	ChainID      string              `json:"chain_id"`
	BranchID     string              `json:"branch_id"`
	TableID      *string             `json:"table_id,omitempty"`
	OrderType    string              `json:"order_type"`
	CustomerName string              `json:"customer_name,omitempty"`
	Items        []orderItemRequest  `json:"items"`
}

type orderItemRequest struct {
	MenuItemID string  `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
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

	if req.ChainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id is required")
		return
	}
	if req.BranchID == "" {
		response.Error(w, http.StatusBadRequest, "branch_id is required")
		return
	}

	var items []domain.OrderItem
	for _, ir := range req.Items {
		items = append(items, domain.OrderItem{
			MenuItemID: ir.MenuItemID,
			Quantity:   ir.Quantity,
			UnitPrice:  ir.UnitPrice,
		})
	}

	order, err := h.svc.Create(req.ChainID, req.BranchID, req.OrderType, req.TableID, nil, nil, req.CustomerName, items)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, order)
}

type statusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id required")
		return
	}

	var req statusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	order, err := h.svc.UpdateStatus(id, req.Status)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, order)
}

type payRequest struct {
	Method string `json:"method"`
}

func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id required")
		return
	}

	var req payRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	order, err := h.svc.Pay(id, req.Method)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, order)
}

func parsePagination(r *http.Request) (page, limit int, err error) {
	page = 1
	limit = 20

	if p := r.URL.Query().Get("page"); p != "" {
		page, err = strconv.Atoi(p)
		if err != nil || page < 1 {
			return 0, 0, err
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil || limit < 1 {
			return 0, 0, err
		}
		if limit > 100 {
			limit = 100
		}
	}
	return page, limit, nil
}

func (h *Handler) ListByBranch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	branchID := r.URL.Query().Get("branch_id")
	if branchID == "" {
		response.Error(w, http.StatusBadRequest, "branch_id required")
		return
	}

	page, limit, err := parsePagination(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid page or limit")
		return
	}

	orders, total, err := h.svc.ListByBranchID(branchID, page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, domain.PaginatedResponse[domain.Order]{
		Items: orders,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *Handler) ListByChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id required")
		return
	}

	page, limit, err := parsePagination(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid page or limit")
		return
	}

	orders, total, err := h.svc.ListByChainID(chainID, page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, domain.PaginatedResponse[domain.Order]{
		Items: orders,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}
