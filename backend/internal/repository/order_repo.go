package repository

import (
	"database/sql"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(order *domain.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow(
		`INSERT INTO orders (chain_id, branch_id, table_id, customer_id, guest_id, customer_name, order_type, status, total_amount)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, created_at, updated_at`,
		order.ChainID, order.BranchID, order.TableID, order.CustomerID, order.GuestID,
		order.CustomerName, order.OrderType, "pending", order.TotalAmount,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		order.Items[i].Subtotal = float64(order.Items[i].Quantity) * order.Items[i].UnitPrice
		err = tx.QueryRow(
			`INSERT INTO order_items (order_id, menu_item_id, variant_id, quantity, unit_price, subtotal)
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			order.ID, order.Items[i].MenuItemID, order.Items[i].VariantID,
			order.Items[i].Quantity, order.Items[i].UnitPrice, order.Items[i].Subtotal,
		).Scan(&order.Items[i].ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepo) GetByID(id string) (*domain.Order, error) {
	o := &domain.Order{}
	err := r.db.QueryRow(
		`SELECT id, chain_id, branch_id, table_id, customer_id,
		        COALESCE(guest_id, '') as guest_id,
		        COALESCE(customer_name, '') as customer_name,
		        order_type, status, payment_status, payment_method,
		        total_amount, created_at, updated_at
		 FROM orders WHERE id = $1`, id,
	).Scan(&o.ID, &o.ChainID, &o.BranchID, &o.TableID, &o.CustomerID, &o.GuestID, &o.CustomerName,
		&o.OrderType, &o.Status, &o.PaymentStatus, &o.PaymentMethod, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *OrderRepo) ListByBranchID(branchID string) ([]domain.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, branch_id, table_id, customer_id,
		        COALESCE(guest_id, '') as guest_id,
		        COALESCE(customer_name, '') as customer_name,
		        order_type, status, payment_status, payment_method,
		        total_amount, created_at, updated_at
		 FROM orders WHERE branch_id = $1 ORDER BY created_at DESC`, branchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.ChainID, &o.BranchID, &o.TableID, &o.CustomerID, &o.GuestID, &o.CustomerName,
			&o.OrderType, &o.Status, &o.PaymentStatus, &o.PaymentMethod, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepo) ListByChainID(chainID string) ([]domain.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, chain_id, branch_id, table_id, customer_id,
		        COALESCE(guest_id, '') as guest_id,
		        COALESCE(customer_name, '') as customer_name,
		        order_type, status, payment_status, payment_method,
		        total_amount, created_at, updated_at
		 FROM orders WHERE chain_id = $1 ORDER BY created_at DESC`, chainID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.ChainID, &o.BranchID, &o.TableID, &o.CustomerID, &o.GuestID, &o.CustomerName,
			&o.OrderType, &o.Status, &o.PaymentStatus, &o.PaymentMethod, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepo) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *OrderRepo) UpdatePayment(id, status, method string) error {
	_, err := r.db.Exec(
		`UPDATE orders SET payment_status=$1, payment_method=$2, updated_at=NOW() WHERE id=$3`,
		status, method, id,
	)
	return err
}
