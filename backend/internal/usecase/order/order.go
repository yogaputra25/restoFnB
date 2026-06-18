package order

import (
	"fmt"

	"github.com/resto-fnb/backend/internal/domain"
)

type OrderRepo interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	ListByBranchID(branchID string) ([]domain.Order, error)
	ListByChainID(chainID string) ([]domain.Order, error)
	UpdateStatus(id, status string) error
	UpdatePayment(id, status, method string) error
}

var validTransitions = map[string][]string{
	"pending":    {"confirmed", "cancelled"},
	"confirmed":  {"preparing", "cancelled"},
	"preparing":  {"ready", "cancelled"},
	"ready":      {"completed"},
	"completed":  {},
	"cancelled":  {},
}

type Usecase struct {
	repo OrderRepo
}

func NewUsecase(repo OrderRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) Create(chainID, branchID, orderType string, tableID, customerID, guestID *string, customerName string, items []domain.OrderItem) (*domain.Order, error) {
	var total float64
	for i := range items {
		items[i].Subtotal = float64(items[i].Quantity) * items[i].UnitPrice
		total += items[i].Subtotal
	}

	var guestIDStr string
	if guestID != nil {
		guestIDStr = *guestID
	}

	order := &domain.Order{
		ChainID:      chainID,
		BranchID:     branchID,
		TableID:      tableID,
		CustomerID:   customerID,
		GuestID:      guestIDStr,
		CustomerName: customerName,
		OrderType:    orderType,
		Status:       "pending",
		TotalAmount:  total,
		Items:        items,
	}

	if err := uc.repo.Create(order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}

func (uc *Usecase) UpdateStatus(id, status string) (*domain.Order, error) {
	order, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	allowed, ok := validTransitions[order.Status]
	if !ok {
		return nil, fmt.Errorf("invalid current status: %s", order.Status)
	}

	valid := false
	for _, s := range allowed {
		if s == status {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("cannot transition from %s to %s", order.Status, status)
	}

	if err := uc.repo.UpdateStatus(id, status); err != nil {
		return nil, err
	}

	order.Status = status
	return order, nil
}

func (uc *Usecase) Pay(id, method string) (*domain.Order, error) {
	order, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	if order.PaymentStatus == "paid" {
		return nil, fmt.Errorf("order already paid")
	}

	if err := uc.repo.UpdatePayment(id, "paid", method); err != nil {
		return nil, err
	}

	order.PaymentStatus = "paid"
	order.PaymentMethod = &method
	return order, nil
}

func (uc *Usecase) ListByBranchID(branchID string) ([]domain.Order, error) {
	return uc.repo.ListByBranchID(branchID)
}

func (uc *Usecase) ListByChainID(chainID string) ([]domain.Order, error) {
	return uc.repo.ListByChainID(chainID)
}
