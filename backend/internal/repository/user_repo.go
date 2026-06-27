package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return r.db.QueryRow(
		`INSERT INTO users (chain_id, branch_id, email, password_hash, name, role, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		user.ChainID, user.BranchID, user.Email, user.PasswordHash, user.Name, user.Role, true,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepo) GetByEmail(chainID, email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(
		`SELECT id, chain_id, branch_id, email, password_hash, name, role, is_active, created_at, updated_at
		 FROM users WHERE chain_id = $1 AND email = $2`,
		chainID, email,
	).Scan(&u.ID, &u.ChainID, &u.BranchID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetByEmailOnly(email string) (*domain.User, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, branch_id, email, password_hash, name, role, is_active, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.ChainID, &u.BranchID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	if len(users) > 1 {
		return nil, fmt.Errorf("multiple accounts found for this email")
	}
	return &users[0], nil
}

func (r *UserRepo) GetByID(id string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(
		`SELECT id, chain_id, branch_id, email, password_hash, name, role, is_active, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.ChainID, &u.BranchID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) ListByChainID(chainID string) ([]domain.User, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, branch_id, email, name, role, is_active, created_at, updated_at
		 FROM users WHERE chain_id = $1 ORDER BY name`,
		chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.ChainID, &u.BranchID, &u.Email, &u.Name, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Deactivate(id string) error {
	_, err := r.db.Exec(`UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *UserRepo) Activate(id string) error {
	_, err := r.db.Exec(`UPDATE users SET is_active = true, updated_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *UserRepo) SaveRefreshToken(token *domain.RefreshToken) error {
	return r.db.QueryRow(
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id`,
		token.UserID, token.TokenHash, token.ExpiresAt,
	).Scan(&token.ID)
}
