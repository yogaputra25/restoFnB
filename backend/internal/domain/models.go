package domain

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCashier  Role = "cashier"
	RoleCustomer Role = "customer"
)

type User struct {
	ID           string    `json:"id"`
	ChainID      string    `json:"chain_id"`
	BranchID     *string   `json:"branch_id,omitempty"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         Role      `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
