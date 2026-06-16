package chain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/resto-fnb/backend/internal/domain"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

type Repo interface {
	Create(chain *domain.Chain) error
	GetBySlug(slug string) (*domain.Chain, error)
	GetByID(id string) (*domain.Chain, error)
}

type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) Register(name, slug, contactInfo string) (*domain.Chain, error) {
	slug = strings.ToLower(strings.TrimSpace(slug))
	if !slugRegex.MatchString(slug) {
		return nil, fmt.Errorf("slug must be lowercase alphanumeric with hyphens")
	}

	existing, _ := uc.repo.GetBySlug(slug)
	if existing != nil {
		return nil, fmt.Errorf("chain with slug %q already exists", slug)
	}

	chain := &domain.Chain{
		Name:        name,
		Slug:        slug,
		ContactInfo: contactInfo,
	}
	if err := uc.repo.Create(chain); err != nil {
		return nil, fmt.Errorf("create chain: %w", err)
	}

	return chain, nil
}

func (uc *Usecase) GetBySlug(slug string) (*domain.Chain, error) {
	return uc.repo.GetBySlug(slug)
}
