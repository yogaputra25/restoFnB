package chain

import (
	"fmt"

	"github.com/resto-fnb/backend/internal/domain"
)

type ThemeRepo interface {
	Upsert(theme *domain.Theme) error
	GetByChainID(chainID string) (*domain.Theme, error)
}

type ThemeUsecase struct {
	themeRepo ThemeRepo
	chainRepo Repo
}

func NewThemeUsecase(themeRepo ThemeRepo, chainRepo Repo) *ThemeUsecase {
	return &ThemeUsecase{themeRepo: themeRepo, chainRepo: chainRepo}
}

func (uc *ThemeUsecase) Upsert(chainID, primaryColor, accentColor, logoURL string) (*domain.Theme, error) {
	if _, err := uc.chainRepo.GetByID(chainID); err != nil {
		return nil, fmt.Errorf("chain not found: %w", err)
	}

	theme := &domain.Theme{
		ChainID:      chainID,
		PrimaryColor: primaryColor,
		AccentColor:  accentColor,
		LogoURL:      logoURL,
	}
	if err := uc.themeRepo.Upsert(theme); err != nil {
		return nil, fmt.Errorf("upsert theme: %w", err)
	}

	return theme, nil
}

func (uc *ThemeUsecase) GetByChainID(chainID string) (*domain.Theme, error) {
	return uc.themeRepo.GetByChainID(chainID)
}
