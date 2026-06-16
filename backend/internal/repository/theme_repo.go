package repository

import (
	"database/sql"

	"github.com/resto-fnb/backend/internal/domain"
)

type ThemeRepo struct {
	db *sql.DB
}

func NewThemeRepo(db *sql.DB) *ThemeRepo {
	return &ThemeRepo{db: db}
}

func (r *ThemeRepo) Upsert(theme *domain.Theme) error {
	return r.db.QueryRow(
		`INSERT INTO themes (chain_id, primary_color, accent_color, logo_url)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (chain_id) DO UPDATE SET
		   primary_color = EXCLUDED.primary_color,
		   accent_color = EXCLUDED.accent_color,
		   logo_url = EXCLUDED.logo_url,
		   updated_at = NOW()
		 RETURNING id`,
		theme.ChainID, theme.PrimaryColor, theme.AccentColor, theme.LogoURL,
	).Scan(&theme.ID)
}

func (r *ThemeRepo) GetByChainID(chainID string) (*domain.Theme, error) {
	t := &domain.Theme{}
	err := r.db.QueryRow(
		`SELECT id, chain_id, primary_color, accent_color, logo_url
		 FROM themes WHERE chain_id = $1`,
		chainID,
	).Scan(&t.ID, &t.ChainID, &t.PrimaryColor, &t.AccentColor, &t.LogoURL)
	if err != nil {
		return nil, err
	}
	return t, nil
}
