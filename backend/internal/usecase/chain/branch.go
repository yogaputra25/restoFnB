package chain

import (
	"fmt"

	"github.com/resto-fnb/backend/internal/domain"
)

type BranchRepo interface {
	Create(branch *domain.Branch) error
	GetByChainID(chainID string) ([]domain.Branch, error)
	GetByID(id string) (*domain.Branch, error)
}

type BranchUsecase struct {
	branchRepo BranchRepo
	chainRepo  Repo
}

func NewBranchUsecase(branchRepo BranchRepo, chainRepo Repo) *BranchUsecase {
	return &BranchUsecase{branchRepo: branchRepo, chainRepo: chainRepo}
}

func (uc *BranchUsecase) Create(chainID, name, address, contactInfo string, parentID *string) (*domain.Branch, error) {
	chain, err := uc.chainRepo.GetByID(chainID)
	if err != nil {
		return nil, fmt.Errorf("chain not found: %w", err)
	}
	_ = chain

	branch := &domain.Branch{
		ChainID:     chainID,
		Name:        name,
		Address:     address,
		ContactInfo: contactInfo,
		ParentID:    parentID,
	}
	if err := uc.branchRepo.Create(branch); err != nil {
		return nil, fmt.Errorf("create branch: %w", err)
	}

	return branch, nil
}

func (uc *BranchUsecase) ListByChainID(chainID string) ([]domain.Branch, error) {
	return uc.branchRepo.GetByChainID(chainID)
}
