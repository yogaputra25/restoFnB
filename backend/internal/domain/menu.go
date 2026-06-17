package domain

import "time"

type Category struct {
	ID           string    `json:"id"`
	ChainID      string    `json:"chain_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MenuItem struct {
	ID          string    `json:"id"`
	ChainID     string    `json:"chain_id"`
	CategoryID  string    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ItemVariant struct {
	ID              string  `json:"id"`
	MenuItemID      string  `json:"menu_item_id"`
	Name            string  `json:"name"`
	PriceAdjustment float64 `json:"price_adjustment"`
}

type BranchAvailability struct {
	MenuItemID  string `json:"menu_item_id"`
	BranchID    string `json:"branch_id"`
	IsAvailable bool   `json:"is_available"`
}

type CategoryStore interface {
	Create(cat *Category) error
	GetByID(id string) (*Category, error)
	ListByChainID(chainID string) ([]Category, error)
	Update(cat *Category) error
	Delete(id string) error
	UpdateOrder(id string, order int) error
}

type MenuItemStore interface {
	Create(item *MenuItem) error
	GetByID(id string) (*MenuItem, error)
	ListByCategoryID(categoryID string) ([]MenuItem, error)
	ListByChainID(chainID string) ([]MenuItem, error)
	Update(item *MenuItem) error
	Delete(id string) error
}

type ItemVariantStore interface {
	Create(variant *ItemVariant) error
	ListByMenuItemID(menuItemID string) ([]ItemVariant, error)
	Delete(id string) error
}

type CarouselSlide struct {
	ID           string     `json:"id"`
	ChainID      string     `json:"chain_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description,omitempty"`
	ImageURL     string     `json:"image_url,omitempty"`
	BgColor      string     `json:"bg_color,omitempty"`
	DisplayOrder int        `json:"display_order"`
	IsActive     bool       `json:"is_active"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CarouselSlideStore interface {
	Create(s *CarouselSlide) error
	ListByChainID(chainID string) ([]CarouselSlide, error)
	ListPublic(chainID string) ([]CarouselSlide, error)
	Update(s *CarouselSlide) error
	SoftDelete(id string) error
}

type BranchAvailabilityStore interface {
	Upsert(ba *BranchAvailability) error
	GetByBranchID(branchID string) ([]BranchAvailability, error)
}
