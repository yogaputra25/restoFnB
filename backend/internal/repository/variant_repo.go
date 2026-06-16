package repository

import (
	"database/sql"

	"github.com/resto-fnb/backend/internal/domain"
)

type VariantRepo struct {
	db *sql.DB
}

func NewVariantRepo(db *sql.DB) *VariantRepo {
	return &VariantRepo{db: db}
}

func (r *VariantRepo) Create(v *domain.ItemVariant) error {
	return r.db.QueryRow(
		`INSERT INTO item_variants (menu_item_id, name, price_adjustment)
		 VALUES ($1, $2, $3) RETURNING id`,
		v.MenuItemID, v.Name, v.PriceAdjustment,
	).Scan(&v.ID)
}

func (r *VariantRepo) ListByMenuItemID(menuItemID string) ([]domain.ItemVariant, error) {
	rows, err := r.db.Query(
		`SELECT id, menu_item_id, name, price_adjustment
		 FROM item_variants WHERE menu_item_id = $1 ORDER BY name`, menuItemID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []domain.ItemVariant
	for rows.Next() {
		var v domain.ItemVariant
		if err := rows.Scan(&v.ID, &v.MenuItemID, &v.Name, &v.PriceAdjustment); err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}
	return variants, nil
}

func (r *VariantRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM item_variants WHERE id = $1`, id)
	return err
}
