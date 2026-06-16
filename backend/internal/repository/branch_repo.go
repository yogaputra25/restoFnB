package repository

import (
	"database/sql"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type BranchRepo struct {
	db *sql.DB
}

func NewBranchRepo(db *sql.DB) *BranchRepo {
	return &BranchRepo{db: db}
}

func (r *BranchRepo) Create(branch *domain.Branch) error {
	branch.CreatedAt = time.Now()
	branch.UpdatedAt = time.Now()

	return r.db.QueryRow(
		`INSERT INTO branches (chain_id, name, address, contact_info, parent_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		branch.ChainID, branch.Name, branch.Address, branch.ContactInfo, branch.ParentID,
	).Scan(&branch.ID, &branch.CreatedAt, &branch.UpdatedAt)
}

func (r *BranchRepo) GetByChainID(chainID string) ([]domain.Branch, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, name, address, contact_info, parent_id, created_at, updated_at
		 FROM branches WHERE chain_id = $1 ORDER BY name`,
		chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var b domain.Branch
		if err := rows.Scan(&b.ID, &b.ChainID, &b.Name, &b.Address, &b.ContactInfo, &b.ParentID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}

func (r *BranchRepo) GetByID(id string) (*domain.Branch, error) {
	b := &domain.Branch{}
	err := r.db.QueryRow(
		`SELECT id, chain_id, name, address, contact_info, parent_id, created_at, updated_at
		 FROM branches WHERE id = $1`,
		id,
	).Scan(&b.ID, &b.ChainID, &b.Name, &b.Address, &b.ContactInfo, &b.ParentID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}
