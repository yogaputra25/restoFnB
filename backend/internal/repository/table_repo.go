package repository

import (
	"database/sql"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type TableRepo struct {
	db *sql.DB
}

func NewTableRepo(db *sql.DB) *TableRepo {
	return &TableRepo{db: db}
}

func (r *TableRepo) Create(t *domain.Table) error {
	t.CreatedAt = time.Now()
	return r.db.QueryRow(
		`INSERT INTO tables (branch_id, label, qr_code_url) VALUES ($1, $2, $3) RETURNING id, created_at`,
		t.BranchID, t.Label, t.QrCodeURL,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TableRepo) ListByBranchID(branchID string) ([]domain.Table, error) {
	rows, err := r.db.Query(
		`SELECT id, branch_id, label, qr_code_url, created_at FROM tables WHERE branch_id = $1 ORDER BY label`,
		branchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []domain.Table
	for rows.Next() {
		var t domain.Table
		if err := rows.Scan(&t.ID, &t.BranchID, &t.Label, &t.QrCodeURL, &t.CreatedAt); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

func (r *TableRepo) GetByID(id string) (*domain.Table, error) {
	t := &domain.Table{}
	err := r.db.QueryRow(
		`SELECT id, branch_id, label, qr_code_url, created_at FROM tables WHERE id = $1`, id,
	).Scan(&t.ID, &t.BranchID, &t.Label, &t.QrCodeURL, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TableRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM tables WHERE id = $1`, id)
	return err
}
