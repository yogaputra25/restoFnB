package usecase

import (
	"github.com/resto-fnb/backend/internal/domain"
)

type CarouselRepo interface {
	Create(s *domain.CarouselSlide) error
	ListByChainID(chainID string) ([]domain.CarouselSlide, error)
	ListPublic(chainID string) ([]domain.CarouselSlide, error)
	Update(s *domain.CarouselSlide) error
	SoftDelete(id string) error
}

type CarouselUsecase struct {
	repo CarouselRepo
}

func NewCarouselUsecase(repo CarouselRepo) *CarouselUsecase {
	return &CarouselUsecase{repo: repo}
}

func (uc *CarouselUsecase) Create(chainID, title, description, imageURL, bgColor string, displayOrder int, isActive bool) (*domain.CarouselSlide, error) {
	s := &domain.CarouselSlide{
		ChainID:      chainID,
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
		BgColor:      bgColor,
		DisplayOrder: displayOrder,
		IsActive:     isActive,
	}
	if err := uc.repo.Create(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *CarouselUsecase) ListByChainID(chainID string) ([]domain.CarouselSlide, error) {
	return uc.repo.ListByChainID(chainID)
}

func (uc *CarouselUsecase) ListPublic(chainID string) ([]domain.CarouselSlide, error) {
	return uc.repo.ListPublic(chainID)
}

func (uc *CarouselUsecase) Update(id, title, description, imageURL, bgColor string, displayOrder int, isActive bool) (*domain.CarouselSlide, error) {
	s := &domain.CarouselSlide{
		ID:           id,
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
		BgColor:      bgColor,
		DisplayOrder: displayOrder,
		IsActive:     isActive,
	}
	if err := uc.repo.Update(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *CarouselUsecase) SoftDelete(id string) error {
	return uc.repo.SoftDelete(id)
}
