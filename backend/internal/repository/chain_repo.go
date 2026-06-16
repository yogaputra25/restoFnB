package repository

import (
	"database/sql"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type ChainRepo struct {
	db *sql.DB
}

func NewChainRepo(db *sql.DB) *ChainRepo {
	return &ChainRepo{db: db}
}

func (r *ChainRepo) Create(chain *domain.Chain) error {
	chain.CreatedAt = time.Now()
	chain.UpdatedAt = time.Now()

	return r.db.QueryRow(
		`INSERT INTO chains (name, slug, contact_info) VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		chain.Name, chain.Slug, chain.ContactInfo,
	).Scan(&chain.ID, &chain.CreatedAt, &chain.UpdatedAt)
}

func (r *ChainRepo) GetBySlug(slug string) (*domain.Chain, error) {
	c := &domain.Chain{}
	err := r.db.QueryRow(
		`SELECT id, name, slug, contact_info, created_at, updated_at FROM chains WHERE slug = $1`,
		slug,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.ContactInfo, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChainRepo) GetByID(id string) (*domain.Chain, error) {
	c := &domain.Chain{}
	err := r.db.QueryRow(
		`SELECT id, name, slug, contact_info, created_at, updated_at FROM chains WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.ContactInfo, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}
