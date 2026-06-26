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

func (r *OrderRepo) ListByBranchID(branchID string, page, limit int) ([]domain.Order, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE branch_id = $1`, branchID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT id, chain_id, branch_id, table_id, customer_id,
		        COALESCE(guest_id, '') as guest_id,
		        COALESCE(customer_name, '') as customer_name,
		        order_type, status, payment_status, payment_method,
		        total_amount, created_at, updated_at
		 FROM orders WHERE branch_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, branchID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.ChainID, &o.BranchID, &o.TableID, &o.CustomerID, &o.GuestID, &o.CustomerName,
			&o.OrderType, &o.Status, &o.PaymentStatus, &o.PaymentMethod, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *OrderRepo) ListByChainID(chainID string, page, limit int) ([]domain.Order, int, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE chain_id = $1`, chainID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT id, chain_id, branch_id, table_id, customer_id,
		        COALESCE(guest_id, '') as guest_id,
		        COALESCE(customer_name, '') as customer_name,
		        order_type, status, payment_status, payment_method,
		        total_amount, created_at, updated_at
		 FROM orders WHERE chain_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, chainID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.ChainID, &o.BranchID, &o.TableID, &o.CustomerID, &o.GuestID, &o.CustomerName,
			&o.OrderType, &o.Status, &o.PaymentStatus, &o.PaymentMethod, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *OrderRepo) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *OrderRepo) GetSalesReport(chainID, period, from, to string) (*domain.SalesReportResponse, error) {
	dateTrunc := "day"
	switch period {
	case "monthly":
		dateTrunc = "month"
	case "yearly":
		dateTrunc = "year"
	}

	report := &domain.SalesReportResponse{
		Period: period,
		From:   from,
		To:     to,
	}

	err := r.db.QueryRow(`
		SELECT
			COALESCE(SUM(o.total_amount), 0),
			COUNT(DISTINCT o.id),
			COALESCE(SUM(oi.quantity), 0)
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id
		WHERE o.chain_id = $1
		  AND o.payment_status = 'paid'
		  AND o.created_at >= $2::timestamp
		  AND o.created_at < ($3::timestamp + INTERVAL '1 day')
	`, chainID, from, to).Scan(&report.Summary.TotalRevenue, &report.Summary.TotalOrders, &report.Summary.TotalItems)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		WITH ranked AS (
			SELECT
				c.id AS category_id,
				c.name AS category_name,
				mi.id AS menu_item_id,
				mi.name,
				SUM(oi.quantity) AS quantity,
				SUM(oi.subtotal) AS revenue,
				ROW_NUMBER() OVER (PARTITION BY c.id ORDER BY SUM(oi.quantity) DESC) AS rn
		 FROM order_items oi
		 JOIN orders o ON o.id = oi.order_id
		 JOIN menu_items mi ON mi.id = oi.menu_item_id
		 JOIN categories c ON c.id = mi.category_id
		 WHERE o.chain_id = $1
		   AND o.payment_status = 'paid'
		   AND o.created_at >= $2::timestamp
		   AND o.created_at < ($3::timestamp + INTERVAL '1 day')
		 GROUP BY c.id, c.name, mi.id, mi.name
		)
		SELECT category_id, category_name, menu_item_id, name, quantity, revenue
		FROM ranked WHERE rn <= 10
		ORDER BY category_name, quantity DESC
	`, chainID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryMap := make(map[string]*domain.SalesReportCategory)
	var categoryOrder []string
	for rows.Next() {
		var catID, catName, itemID, itemName string
		var qty int
		var rev float64
		if err := rows.Scan(&catID, &catName, &itemID, &itemName, &qty, &rev); err != nil {
			return nil, err
		}
		cat, ok := categoryMap[catID]
		if !ok {
			cat = &domain.SalesReportCategory{
				CategoryID:   catID,
				CategoryName: catName,
				Items:        []domain.SalesReportItem{},
			}
			categoryMap[catID] = cat
			categoryOrder = append(categoryOrder, catID)
		}
		cat.Items = append(cat.Items, domain.SalesReportItem{
			MenuItemID: itemID,
			Name:       itemName,
			Quantity:   qty,
			Revenue:    rev,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, id := range categoryOrder {
		report.Categories = append(report.Categories, *categoryMap[id])
	}

	bRows, err := r.db.Query(`
		SELECT
			to_char(date_trunc($4, created_at), 'YYYY-MM-DD') AS date,
			COALESCE(SUM(total_amount), 0),
			COUNT(*)
		FROM orders
		WHERE chain_id = $1
		  AND payment_status = 'paid'
		  AND created_at >= $2::timestamp
		  AND created_at < ($3::timestamp + INTERVAL '1 day')
		GROUP BY date_trunc($4, created_at)
		ORDER BY date
	`, chainID, from, to, dateTrunc)
	if err != nil {
		return nil, err
	}
	defer bRows.Close()

	for bRows.Next() {
		var d domain.SalesReportDaily
		if err := bRows.Scan(&d.Date, &d.Revenue, &d.Orders); err != nil {
			return nil, err
		}
		report.DailyBreakdown = append(report.DailyBreakdown, d)
	}
	if err := bRows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *OrderRepo) UpdatePayment(id, status, method string) error {
	_, err := r.db.Exec(
		`UPDATE orders SET payment_status=$1, payment_method=$2, updated_at=NOW() WHERE id=$3`,
		status, method, id,
	)
	return err
}
