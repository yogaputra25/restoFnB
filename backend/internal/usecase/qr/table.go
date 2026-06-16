package qr

import (
	"fmt"

	"github.com/resto-fnb/backend/internal/domain"
)

type TableRepo interface {
	Create(t *domain.Table) error
	ListByBranchID(branchID string) ([]domain.Table, error)
	GetByID(id string) (*domain.Table, error)
	Delete(id string) error
}

type TableUsecase struct {
	repo       TableRepo
	frontendURL string
}

func NewTableUsecase(repo TableRepo, frontendURL string) *TableUsecase {
	return &TableUsecase{repo: repo, frontendURL: frontendURL}
}

func (uc *TableUsecase) Create(branchID, label string) (*domain.Table, error) {
	t := &domain.Table{
		BranchID: branchID,
		Label:    label,
	}
	if err := uc.repo.Create(t); err != nil {
		return nil, err
	}

	t.QrCodeURL = fmt.Sprintf("%s/order?branch=%s&table=%s", uc.frontendURL, branchID, t.ID)

	return t, nil
}

func (uc *TableUsecase) ListByBranchID(branchID string) ([]domain.Table, error) {
	return uc.repo.ListByBranchID(branchID)
}

func (uc *TableUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}
