package repository

import (
	"database/sql"

	"github.com/resto-fnb/backend/internal/domain"
)

type AvailabilityRepo struct {
	db *sql.DB
}

func NewAvailabilityRepo(db *sql.DB) *AvailabilityRepo {
	return &AvailabilityRepo{db: db}
}

func (r *AvailabilityRepo) Upsert(ba *domain.BranchAvailability) error {
	_, err := r.db.Exec(
		`INSERT INTO branch_availability (menu_item_id, branch_id, is_available)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (menu_item_id, branch_id) DO UPDATE SET is_available = EXCLUDED.is_available`,
		ba.MenuItemID, ba.BranchID, ba.IsAvailable,
	)
	return err
}

func (r *AvailabilityRepo) GetByBranchID(branchID string) ([]domain.BranchAvailability, error) {
	rows, err := r.db.Query(
		`SELECT menu_item_id, branch_id, is_available
		 FROM branch_availability WHERE branch_id = $1`, branchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.BranchAvailability
	for rows.Next() {
		var ba domain.BranchAvailability
		if err := rows.Scan(&ba.MenuItemID, &ba.BranchID, &ba.IsAvailable); err != nil {
			return nil, err
		}
		items = append(items, ba)
	}
	return items, nil
}
