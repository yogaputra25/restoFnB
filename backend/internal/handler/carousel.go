package handler

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/pkg/response"
)

type CarouselHandler struct {
	create func(chainID, title, description, imageURL, bgColor string, displayOrder int, isActive bool) (*domain.CarouselSlide, error)
	list   func(chainID string) ([]domain.CarouselSlide, error)
	public func(chainID string) ([]domain.CarouselSlide, error)
	update func(id, title, description, imageURL, bgColor string, displayOrder int, isActive bool) (*domain.CarouselSlide, error)
	delete func(id string) error
}

func NewCarouselHandler(
	create func(chainID, title, description, imageURL, bgColor string, displayOrder int, isActive bool) (*domain.CarouselSlide, error),
	list func(chainID string) ([]domain.CarouselSlide, error),
	public func(chainID string) ([]domain.CarouselSlide, error),
	update func(id, title, description, imageURL, bgColor string, displayOrder int, isActive bool) (*domain.CarouselSlide, error),
	delete func(id string) error,
) *CarouselHandler {
	return &CarouselHandler{create: create, list: list, public: public, update: update, delete: delete}
}

type carouselRequest struct {
	ChainID      string `json:"chain_id"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	BgColor      string `json:"bg_color,omitempty"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

type idRequest struct {
	ID string `json:"id"`
}

type carouselUpdateRequest struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	BgColor      string `json:"bg_color,omitempty"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

func (h *CarouselHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req carouselRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	s, err := h.create(req.ChainID, req.Title, req.Description, req.ImageURL, req.BgColor, req.DisplayOrder, req.IsActive)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, s)
}

func (h *CarouselHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id required")
		return
	}
	slides, err := h.list(chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, slides)
}

func (h *CarouselHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id required")
		return
	}
	slides, err := h.public(chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, slides)
}

func (h *CarouselHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req carouselUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	s, err := h.update(req.ID, req.Title, req.Description, req.ImageURL, req.BgColor, req.DisplayOrder, req.IsActive)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, s)
}

func (h *CarouselHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	response.JSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
