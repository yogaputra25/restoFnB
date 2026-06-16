package domain

import "time"

type Table struct {
	ID        string    `json:"id"`
	BranchID  string    `json:"branch_id"`
	Label     string    `json:"label"`
	QrCodeURL string    `json:"qr_code_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TableStore interface {
	Create(t *Table) error
	ListByBranchID(branchID string) ([]Table, error)
	GetByID(id string) (*Table, error)
	Delete(id string) error
}
