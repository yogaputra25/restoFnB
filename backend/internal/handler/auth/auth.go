package auth

import (
	"encoding/json"
	"net/http"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/internal/middleware"
	"github.com/resto-fnb/backend/pkg/response"
)

type RegisterRequest struct {
	ChainID  string       `json:"chain_id"`
	BranchID *string      `json:"branch_id,omitempty"`
	Email    string       `json:"email"`
	Password string       `json:"password"`
	Name     string       `json:"name"`
	Role     domain.Role  `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Register(chainID string, branchID *string, email, password, name string, role domain.Role) (*domain.User, error)
	Login(email, password string) (*domain.AuthTokens, error)
	Refresh(refreshToken string) (*domain.AuthTokens, error)
	GetProfile(userID string) (*domain.User, error)
	ListStaff(chainID string) ([]domain.User, error)
	DeactivateStaff(id string) error
}

type Handler struct {
	svc AuthService
}

func NewHandler(svc AuthService) *Handler {
	return &Handler{svc: svc}
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

	user, err := h.svc.Register(req.ChainID, req.BranchID, req.Email, req.Password, req.Name, req.Role)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.Created(w, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	tokens, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, tokens)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	tokens, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, tokens)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, profile)
}

func (h *Handler) ListStaff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		response.Error(w, http.StatusBadRequest, "chain_id required")
		return
	}

	users, err := h.svc.ListStaff(chainID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, users)
}

func (h *Handler) DeactivateStaff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.svc.DeactivateStaff(id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, map[string]string{"status": "deactivated"})
}
