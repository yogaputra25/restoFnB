package menu

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type CategoryHandler struct {
	create func(chainID, name, description string, displayOrder int) (*domain.Category, error)
	list   func(chainID string) ([]domain.Category, error)
	update func(id, name, description string, displayOrder int) (*domain.Category, error)
	delete func(id string) error
}

func NewCategoryHandler(
	create func(chainID, name, description string, displayOrder int) (*domain.Category, error),
	list func(chainID string) ([]domain.Category, error),
	update func(id, name, description string, displayOrder int) (*domain.Category, error),
	delete func(id string) error,
) *CategoryHandler {
	return &CategoryHandler{create: create, list: list, update: update, delete: delete}
}

type categoryRequest struct {
	ChainID      string `json:"chain_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	DisplayOrder int    `json:"display_order"`
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	cat, err := h.create(req.ChainID, req.Name, req.Description, req.DisplayOrder)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(w, cat)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id required")
		return
	}
	cats, err := h.list(chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, cats)
}

type MenuItemHandler struct {
	create     func(chainID, categoryID, name, description string, price float64, imageURL string) (*domain.MenuItem, error)
	list       func(chainID string) ([]domain.MenuItem, error)
	listPublic func(chainID string) ([]domain.MenuItem, error)
	update     func(id, categoryID, name, description string, price float64, imageURL string, isAvailable bool) (*domain.MenuItem, error)
	delete     func(id string) error
}

func NewMenuItemHandler(
	create func(chainID, categoryID, name, description string, price float64, imageURL string) (*domain.MenuItem, error),
	list func(chainID string) ([]domain.MenuItem, error),
	listPublic func(chainID string) ([]domain.MenuItem, error),
	update func(id, categoryID, name, description string, price float64, imageURL string, isAvailable bool) (*domain.MenuItem, error),
	delete func(id string) error,
) *MenuItemHandler {
	return &MenuItemHandler{create: create, list: list, listPublic: listPublic, update: update, delete: delete}
}

type menuItemRequest struct {
	ChainID     string  `json:"chain_id"`
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url,omitempty"`
	IsAvailable bool    `json:"is_available"`
}

func (h *MenuItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req menuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	item, err := h.create(req.ChainID, req.CategoryID, req.Name, req.Description, req.Price, req.ImageURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(w, item)
}

func (h *MenuItemHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id required")
		return
	}
	items, err := h.listPublic(chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, items)
}

func (h *MenuItemHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id required")
		return
	}
	items, err := h.list(chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, items)
}

type VariantHandler struct {
	create func(menuItemID, name string, priceAdjustment float64) (*domain.ItemVariant, error)
	list   func(menuItemID string) ([]domain.ItemVariant, error)
	delete func(id string) error
}

func NewVariantHandler(
	create func(menuItemID, name string, priceAdjustment float64) (*domain.ItemVariant, error),
	list func(menuItemID string) ([]domain.ItemVariant, error),
	delete func(id string) error,
) *VariantHandler {
	return &VariantHandler{create: create, list: list, delete: delete}
}

type variantRequest struct {
	MenuItemID      string  `json:"menu_item_id"`
	Name            string  `json:"name"`
	PriceAdjustment float64 `json:"price_adjustment"`
}

func (h *VariantHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req variantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	v, err := h.create(req.MenuItemID, req.Name, req.PriceAdjustment)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(w, v)
}

func (h *VariantHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	menuItemID := r.URL.Query().Get("menu_item_id")
	if menuItemID == "" {
		response.Error(w, http.StatusBadRequest, "menu_item_id required")
		return
	}
	variants, err := h.list(menuItemID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, variants)
}

type idRequest struct {
	ID string `json:"id"`
}

type categoryUpdateRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	DisplayOrder int    `json:"display_order"`
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req idRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.delete(req.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, map[string]string{"message": "deleted"})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req categoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	cat, err := h.update(req.ID, req.Name, req.Description, req.DisplayOrder)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, cat)
}

type menuItemUpdateRequest struct {
	ID          string  `json:"id"`
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url,omitempty"`
	IsAvailable bool    `json:"is_available"`
}

func (h *MenuItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req idRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.delete(req.ID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, map[string]string{"message": "deleted"})
}

func (h *MenuItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req menuItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	item, err := h.update(req.ID, req.CategoryID, req.Name, req.Description, req.Price, req.ImageURL, req.IsAvailable)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, item)
}

type AvailabilityHandler struct {
	set  func(menuItemID, branchID string, isAvailable bool) (*domain.BranchAvailability, error)
	get  func(branchID string) ([]domain.BranchAvailability, error)
}

func NewAvailabilityHandler(
	set func(menuItemID, branchID string, isAvailable bool) (*domain.BranchAvailability, error),
	get func(branchID string) ([]domain.BranchAvailability, error),
) *AvailabilityHandler {
	return &AvailabilityHandler{set: set, get: get}
}

type availabilityRequest struct {
	MenuItemID  string `json:"menu_item_id"`
	BranchID    string `json:"branch_id"`
	IsAvailable bool   `json:"is_available"`
}

func (h *AvailabilityHandler) Set(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req availabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	ba, err := h.set(req.MenuItemID, req.BranchID, req.IsAvailable)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, ba)
}

func (h *AvailabilityHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	branchID := r.URL.Query().Get("branch_id")
	if branchID == "" {
		response.Error(w, http.StatusBadRequest, "branch_id required")
		return
	}
	items, err := h.get(branchID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, items)
}
