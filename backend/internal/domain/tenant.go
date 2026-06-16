package domain

import "time"

type Chain struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	ContactInfo string    `json:"contact_info,omitempty"`
	DBName      string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Branch struct {
	ID          string    `json:"id"`
	ChainID     string    `json:"chain_id"`
	Name        string    `json:"name"`
	Address     string    `json:"address,omitempty"`
	ContactInfo string    `json:"contact_info,omitempty"`
	ParentID    *string   `json:"parent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Theme struct {
	ID           string `json:"id"`
	ChainID      string `json:"chain_id"`
	PrimaryColor string `json:"primary_color"`
	AccentColor  string `json:"accent_color"`
	LogoURL      string `json:"logo_url,omitempty"`
}

type TenantStore interface {
	RegisterChain(chain *Chain) error
	GetChainBySlug(slug string) (*Chain, error)
	ProvisionDatabase(chain *Chain) error
}

type BranchStore interface {
	Create(branch *Branch) error
	GetByChainID(chainID string) ([]Branch, error)
	GetByID(id string) (*Branch, error)
}

type ThemeStore interface {
	Upsert(theme *Theme) error
	GetByChainID(chainID string) (*Theme, error)
}
