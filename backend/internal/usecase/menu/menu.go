package menu

import "github.com/resto-fnb/backend/internal/domain"

type CategoryRepo interface {
	Create(cat *domain.Category) error
	GetByID(id string) (*domain.Category, error)
	ListByChainID(chainID string) ([]domain.Category, error)
	Update(cat *domain.Category) error
	Delete(id string) error
	UpdateOrder(id string, order int) error
}

type MenuItemRepo interface {
	Create(item *domain.MenuItem) error
	GetByID(id string) (*domain.MenuItem, error)
	ListByCategoryID(categoryID string) ([]domain.MenuItem, error)
	ListByChainID(chainID string) ([]domain.MenuItem, error)
	Update(item *domain.MenuItem) error
	Delete(id string) error
}

type VariantRepo interface {
	Create(v *domain.ItemVariant) error
	ListByMenuItemID(menuItemID string) ([]domain.ItemVariant, error)
	Delete(id string) error
}

type AvailabilityRepo interface {
	Upsert(ba *domain.BranchAvailability) error
	GetByBranchID(branchID string) ([]domain.BranchAvailability, error)
}

type CategoryUsecase struct {
	repo CategoryRepo
}

func NewCategoryUsecase(repo CategoryRepo) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (uc *CategoryUsecase) Create(chainID, name, description string, displayOrder int) (*domain.Category, error) {
	cat := &domain.Category{
		ChainID:      chainID,
		Name:         name,
		Description:  description,
		DisplayOrder: displayOrder,
	}
	if err := uc.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (uc *CategoryUsecase) List(chainID string) ([]domain.Category, error) {
	return uc.repo.ListByChainID(chainID)
}

func (uc *CategoryUsecase) Update(id, name, description string, displayOrder int) (*domain.Category, error) {
	cat, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	cat.Name = name
	cat.Description = description
	cat.DisplayOrder = displayOrder
	if err := uc.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (uc *CategoryUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}

type MenuItemUsecase struct {
	repo MenuItemRepo
}

func NewMenuItemUsecase(repo MenuItemRepo) *MenuItemUsecase {
	return &MenuItemUsecase{repo: repo}
}

func (uc *MenuItemUsecase) Create(chainID, categoryID, name, description string, price float64, imageURL string) (*domain.MenuItem, error) {
	item := &domain.MenuItem{
		ChainID:     chainID,
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		Price:       price,
		ImageURL:    imageURL,
		IsAvailable: true,
	}
	if err := uc.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (uc *MenuItemUsecase) ListByChainID(chainID string) ([]domain.MenuItem, error) {
	return uc.repo.ListByChainID(chainID)
}

func (uc *MenuItemUsecase) Update(id, categoryID, name, description string, price float64, imageURL string, isAvailable bool) (*domain.MenuItem, error) {
	item, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	item.CategoryID = categoryID
	item.Name = name
	item.Description = description
	item.Price = price
	item.ImageURL = imageURL
	item.IsAvailable = isAvailable
	if err := uc.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (uc *MenuItemUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}

type VariantUsecase struct {
	repo VariantRepo
}

func NewVariantUsecase(repo VariantRepo) *VariantUsecase {
	return &VariantUsecase{repo: repo}
}

func (uc *VariantUsecase) Create(menuItemID, name string, priceAdjustment float64) (*domain.ItemVariant, error) {
	v := &domain.ItemVariant{
		MenuItemID:      menuItemID,
		Name:            name,
		PriceAdjustment: priceAdjustment,
	}
	if err := uc.repo.Create(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (uc *VariantUsecase) ListByMenuItemID(menuItemID string) ([]domain.ItemVariant, error) {
	return uc.repo.ListByMenuItemID(menuItemID)
}

func (uc *VariantUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}

type AvailabilityUsecase struct {
	repo AvailabilityRepo
}

func NewAvailabilityUsecase(repo AvailabilityRepo) *AvailabilityUsecase {
	return &AvailabilityUsecase{repo: repo}
}

func (uc *AvailabilityUsecase) SetAvailability(menuItemID, branchID string, isAvailable bool) (*domain.BranchAvailability, error) {
	ba := &domain.BranchAvailability{
		MenuItemID:  menuItemID,
		BranchID:    branchID,
		IsAvailable: isAvailable,
	}
	if err := uc.repo.Upsert(ba); err != nil {
		return nil, err
	}
	return ba, nil
}

func (uc *AvailabilityUsecase) GetByBranchID(branchID string) ([]domain.BranchAvailability, error) {
	return uc.repo.GetByBranchID(branchID)
}
