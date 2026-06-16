package repository

import (
	"database/sql"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(cat *domain.Category) error {
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()

	return r.db.QueryRow(
		`INSERT INTO categories (chain_id, name, description, display_order)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at, updated_at`,
		cat.ChainID, cat.Name, cat.Description, cat.DisplayOrder,
	).Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt)
}

func (r *CategoryRepo) GetByID(id string) (*domain.Category, error) {
	c := &domain.Category{}
	err := r.db.QueryRow(
		`SELECT id, chain_id, name, description, display_order, created_at, updated_at
		 FROM categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.ChainID, &c.Name, &c.Description, &c.DisplayOrder, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CategoryRepo) ListByChainID(chainID string) ([]domain.Category, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, name, description, display_order, created_at, updated_at
		 FROM categories WHERE chain_id = $1 ORDER BY display_order, name`, chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.ChainID, &c.Name, &c.Description, &c.DisplayOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *CategoryRepo) Update(cat *domain.Category) error {
	cat.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE categories SET name=$1, description=$2, display_order=$3, updated_at=$4 WHERE id=$5`,
		cat.Name, cat.Description, cat.DisplayOrder, cat.UpdatedAt, cat.ID,
	)
	return err
}

func (r *CategoryRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}

func (r *CategoryRepo) UpdateOrder(id string, order int) error {
	_, err := r.db.Exec(`UPDATE categories SET display_order=$1, updated_at=NOW() WHERE id=$2`, order, id)
	return err
}
