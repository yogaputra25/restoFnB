package qr

import (
	"fmt"
	"strings"

	"github.com/resto-fnb/backend/internal/domain"
)

type TableRepo interface {
	Create(t *domain.Table) error
	ListByBranchID(branchID string) ([]domain.Table, error)
	GetByID(id string) (*domain.Table, error)
	GetBranchByID(id string) (*domain.Branch, error)
	Update(id, label string) error
	UpdateQrCodeURL(id, url string) error
	Delete(id string) error
}

func (uc *TableUsecase) qrURL(chainID, branchID, tableID string) string {
	return fmt.Sprintf("%s/menu?chain_id=%s&branch_id=%s&table_id=%s", uc.frontendURL, chainID, branchID, tableID)
}

type TableUsecase struct {
	repo       TableRepo
	frontendURL string
}

func NewTableUsecase(repo TableRepo, frontendURL string) *TableUsecase {
	return &TableUsecase{repo: repo, frontendURL: frontendURL}
}

func (uc *TableUsecase) Create(chainID, branchID, label string) (*domain.Table, error) {
	t := &domain.Table{
		BranchID: branchID,
		Label:    label,
	}
	if err := uc.repo.Create(t); err != nil {
		return nil, err
	}

	t.QrCodeURL = uc.qrURL(chainID, branchID, t.ID)
	if err := uc.repo.UpdateQrCodeURL(t.ID, t.QrCodeURL); err != nil {
		return nil, err
	}

	return t, nil
}

func (uc *TableUsecase) Update(id, label string) (*domain.Table, error) {
	if err := uc.repo.Update(id, label); err != nil {
		return nil, err
	}
	t, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (uc *TableUsecase) ListByBranchID(branchID string) ([]domain.Table, error) {
	tables, err := uc.repo.ListByBranchID(branchID)
	if err != nil {
		return nil, err
	}
	for i := range tables {
		if tables[i].QrCodeURL == "" || !strings.HasPrefix(tables[i].QrCodeURL, uc.frontendURL) {
			branch, err := uc.repo.GetBranchByID(tables[i].BranchID)
			if err != nil {
				continue
			}
			tables[i].QrCodeURL = uc.qrURL(branch.ChainID, tables[i].BranchID, tables[i].ID)
			_ = uc.repo.UpdateQrCodeURL(tables[i].ID, tables[i].QrCodeURL)
		}
	}
	return tables, nil
}

func (uc *TableUsecase) RegenerateQrURL(id string) (*domain.Table, error) {
	t, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	branch, err := uc.repo.GetBranchByID(t.BranchID)
	if err != nil {
		return nil, err
	}
	t.QrCodeURL = uc.qrURL(branch.ChainID, t.BranchID, t.ID)
	if err := uc.repo.UpdateQrCodeURL(t.ID, t.QrCodeURL); err != nil {
		return nil, err
	}
	return t, nil
}

func (uc *TableUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}
