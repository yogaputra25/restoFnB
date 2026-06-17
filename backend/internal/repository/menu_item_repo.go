package repository

import (
	"database/sql"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type MenuItemRepo struct {
	db *sql.DB
}

func NewMenuItemRepo(db *sql.DB) *MenuItemRepo {
	return &MenuItemRepo{db: db}
}

func (r *MenuItemRepo) Create(item *domain.MenuItem) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	return r.db.QueryRow(
		`INSERT INTO menu_items (chain_id, category_id, name, description, price, image_url, is_available)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		item.ChainID, item.CategoryID, item.Name, item.Description, item.Price, item.ImageURL, item.IsAvailable,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *MenuItemRepo) GetByID(id string) (*domain.MenuItem, error) {
	item := &domain.MenuItem{}
	err := r.db.QueryRow(
		`SELECT id, chain_id, category_id, name, description, price, image_url, is_available, created_at, updated_at
		 FROM menu_items WHERE id = $1`, id,
	).Scan(&item.ID, &item.ChainID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.ImageURL, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *MenuItemRepo) ListByCategoryID(categoryID string) ([]domain.MenuItem, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, category_id, name, description, price, image_url, is_available, created_at, updated_at
		 FROM menu_items WHERE category_id = $1 ORDER BY name`, categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var item domain.MenuItem
		if err := rows.Scan(&item.ID, &item.ChainID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.ImageURL, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MenuItemRepo) ListByChainID(chainID string) ([]domain.MenuItem, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, category_id, name, description, price, image_url, is_available, created_at, updated_at
		 FROM menu_items WHERE chain_id = $1 ORDER BY category_id, name`, chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var item domain.MenuItem
		if err := rows.Scan(&item.ID, &item.ChainID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.ImageURL, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MenuItemRepo) ListPublicByChainID(chainID string) ([]domain.MenuItem, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, category_id, name, description, price, image_url, is_available, created_at, updated_at
		 FROM menu_items WHERE chain_id = $1 AND is_available = true ORDER BY category_id, name`, chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var item domain.MenuItem
		if err := rows.Scan(&item.ID, &item.ChainID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.ImageURL, &item.IsAvailable, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MenuItemRepo) Update(item *domain.MenuItem) error {
	item.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE menu_items SET category_id=$1, name=$2, description=$3, price=$4, image_url=$5, is_available=$6, updated_at=$7 WHERE id=$8`,
		item.CategoryID, item.Name, item.Description, item.Price, item.ImageURL, item.IsAvailable, item.UpdatedAt, item.ID,
	)
	return err
}

func (r *MenuItemRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM menu_items WHERE id = $1`, id)
	return err
}
