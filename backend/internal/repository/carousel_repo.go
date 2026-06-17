package repository

import (
	"database/sql"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type CarouselRepo struct {
	db *sql.DB
}

func NewCarouselRepo(db *sql.DB) *CarouselRepo {
	return &CarouselRepo{db: db}
}

func (r *CarouselRepo) Create(s *domain.CarouselSlide) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return r.db.QueryRow(
		`INSERT INTO carousel_slides (chain_id, title, description, image_url, bg_color, display_order, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		s.ChainID, s.Title, s.Description, s.ImageURL, s.BgColor, s.DisplayOrder, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *CarouselRepo) ListByChainID(chainID string) ([]domain.CarouselSlide, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, title, description, image_url, bg_color, display_order, is_active, deleted_at, created_at, updated_at
		 FROM carousel_slides WHERE chain_id = $1 AND deleted_at IS NULL
		 ORDER BY display_order, created_at`, chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slides []domain.CarouselSlide
	for rows.Next() {
		var s domain.CarouselSlide
		if err := rows.Scan(&s.ID, &s.ChainID, &s.Title, &s.Description, &s.ImageURL, &s.BgColor, &s.DisplayOrder, &s.IsActive, &s.DeletedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slides = append(slides, s)
	}
	return slides, nil
}

func (r *CarouselRepo) ListPublic(chainID string) ([]domain.CarouselSlide, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, title, description, image_url, bg_color, display_order, is_active, deleted_at, created_at, updated_at
		 FROM carousel_slides WHERE chain_id = $1 AND is_active = true AND deleted_at IS NULL
		 ORDER BY display_order, created_at`, chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slides []domain.CarouselSlide
	for rows.Next() {
		var s domain.CarouselSlide
		if err := rows.Scan(&s.ID, &s.ChainID, &s.Title, &s.Description, &s.ImageURL, &s.BgColor, &s.DisplayOrder, &s.IsActive, &s.DeletedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slides = append(slides, s)
	}
	return slides, nil
}

func (r *CarouselRepo) Update(s *domain.CarouselSlide) error {
	s.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE carousel_slides SET title=$1, description=$2, image_url=$3, bg_color=$4, display_order=$5, is_active=$6, updated_at=$7 WHERE id=$8 AND deleted_at IS NULL`,
		s.Title, s.Description, s.ImageURL, s.BgColor, s.DisplayOrder, s.IsActive, s.UpdatedAt, s.ID,
	)
	return err
}

func (r *CarouselRepo) SoftDelete(id string) error {
	_, err := r.db.Exec(`UPDATE carousel_slides SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}
