package seed

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/resto-fnb/backend/internal/domain"
	"github.com/resto-fnb/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func Run(db *sql.DB) error {
	log.Println("seeding database...")

	chainRepo := repository.NewChainRepo(db)
	branchRepo := repository.NewBranchRepo(db)
	userRepo := repository.NewUserRepo(db)
	catRepo := repository.NewCategoryRepo(db)
	itemRepo := repository.NewMenuItemRepo(db)
	varRepo := repository.NewVariantRepo(db)
	tableRepo := repository.NewTableRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	var rowCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM chains").Scan(&rowCount); err != nil {
		return fmt.Errorf("seed check chains: %w", err)
	}
	if rowCount > 0 {
		log.Println("database already seeded, skipping")
		return nil
	}

	chain := &domain.Chain{
		Name:        "Seed Chain",
		Slug:        "seed-chain",
		ContactInfo: "admin@seed.com",
	}
	if err := chainRepo.Create(chain); err != nil {
		return fmt.Errorf("seed chain: %w", err)
	}
	log.Printf("created chain: %s (%s)", chain.Name, chain.ID)

	branchA := &domain.Branch{
		ChainID:     chain.ID,
		Name:        "Seed Branch A",
		Address:     "Jl. Makan Enak No. 1",
		ContactInfo: "021-111111",
	}
	if err := branchRepo.Create(branchA); err != nil {
		return fmt.Errorf("seed branch A: %w", err)
	}

	branchB := &domain.Branch{
		ChainID:     chain.ID,
		Name:        "Seed Branch B",
		Address:     "Jl. Makan Enak No. 2",
		ContactInfo: "021-222222",
	}
	if err := branchRepo.Create(branchB); err != nil {
		return fmt.Errorf("seed branch B: %w", err)
	}

	log.Printf("created branches: %s, %s", branchA.Name, branchB.Name)

	seedUsers(userRepo, chain.ID, branchA.ID)

	makananID, minumanID := seedCategories(catRepo, chain.ID)

	nasiGorengID := seedMenuItem(itemRepo, chain.ID, makananID, "Nasi Goreng", "Nasi goreng spesial dengan telur dan ayam", 35000)
	mieGorengID := seedMenuItem(itemRepo, chain.ID, makananID, "Mie Goreng", "Mie goreng dengan sayuran dan udang", 30000)
	esTehID := seedMenuItem(itemRepo, chain.ID, minumanID, "Es Teh Manis", "Teh manis dengan es batu", 5000)
	seedMenuItem(itemRepo, chain.ID, minumanID, "Jus Jeruk", "Jus jeruk segar peras", 12000)

	seedVariant(varRepo, nasiGorengID, "Biasa", 0)
	seedVariant(varRepo, nasiGorengID, "Spesial", 10000)
	seedVariant(varRepo, mieGorengID, "Biasa", 0)
	seedVariant(varRepo, mieGorengID, "Spesial", 8000)

	log.Printf("created categories, menu items, and variants")

	tablesA := seedTables(tableRepo, branchA.ID)
	seedTables(tableRepo, branchB.ID)

	log.Printf("created %d tables", len(tablesA))

	seedOrders(orderRepo, chain.ID, branchA.ID, tablesA, nasiGorengID, esTehID)

	log.Println("seed complete")
	return nil
}

func seedUsers(repo *repository.UserRepo, chainID, branchID string) {
	hash := hashPassword("seed123")

	users := []struct {
		email string
		name  string
		role  domain.Role
	}{
		{"admin@seed.com", "Admin Seed", domain.RoleAdmin},
		{"cashier@seed.com", "Cashier Seed", domain.RoleCashier},
	}

	for _, u := range users {
		if _, err := repo.GetByEmail(chainID, u.email); err == nil {
			continue
		}
		var bID *string
		if u.role == domain.RoleCashier {
			bID = &branchID
		}
		user := &domain.User{
			ChainID:      chainID,
			BranchID:     bID,
			Email:        u.email,
			PasswordHash: hash,
			Name:         u.name,
			Role:         u.role,
			IsActive:     true,
		}
		if err := repo.Create(user); err != nil {
			log.Printf("seed user %s: %v", u.email, err)
		}
	}
}

func seedCategories(repo *repository.CategoryRepo, chainID string) (makananID, minumanID string) {
	cats := []struct {
		name  string
		order int
	}{
		{"Makanan", 1},
		{"Minuman", 2},
	}

	for _, c := range cats {
		cat := &domain.Category{
			ChainID:      chainID,
			Name:         c.name,
			DisplayOrder: c.order,
		}
		if err := repo.Create(cat); err != nil {
			log.Printf("seed category %s: %v", c.name, err)
		}
		if c.name == "Makanan" {
			makananID = cat.ID
		} else {
			minumanID = cat.ID
		}
	}
	return
}

func seedMenuItem(repo *repository.MenuItemRepo, chainID, categoryID, name, desc string, price float64) string {
	item := &domain.MenuItem{
		ChainID:     chainID,
		CategoryID:  categoryID,
		Name:        name,
		Description: desc,
		Price:       price,
		IsAvailable: true,
	}
	if err := repo.Create(item); err != nil {
		log.Printf("seed menu item %s: %v", name, err)
	}
	return item.ID
}

func seedVariant(repo *repository.VariantRepo, menuItemID, name string, adjustment float64) {
	v := &domain.ItemVariant{
		MenuItemID:      menuItemID,
		Name:            name,
		PriceAdjustment: adjustment,
	}
	if err := repo.Create(v); err != nil {
		log.Printf("seed variant %s: %v", name, err)
	}
}

func seedTables(repo *repository.TableRepo, branchID string) []domain.Table {
	var tables []domain.Table
	for i := 1; i <= 2; i++ {
		t := &domain.Table{
			BranchID: branchID,
			Label:    fmt.Sprintf("Meja %d", i),
		}
		if err := repo.Create(t); err != nil {
			log.Printf("seed table Meja %d: %v", i, err)
		}
		tables = append(tables, *t)
	}
	return tables
}

func seedOrders(repo *repository.OrderRepo, chainID, branchID string, tables []domain.Table, menuItemID, drinkID string) {
	if len(tables) < 2 {
		return
	}

	pendingOrder := &domain.Order{
		ChainID:     chainID,
		BranchID:    branchID,
		TableID:     &tables[0].ID,
		OrderType:   "dine-in",
		TotalAmount: 45000,
		CreatedAt:   time.Now(),
		Items: []domain.OrderItem{
			{MenuItemID: menuItemID, Quantity: 1, UnitPrice: 35000, Subtotal: 35000},
			{MenuItemID: drinkID, Quantity: 2, UnitPrice: 5000, Subtotal: 10000},
		},
	}
	if err := repo.Create(pendingOrder); err != nil {
		log.Printf("seed pending order: %v", err)
	}

	completedOrder := &domain.Order{
		ChainID:     chainID,
		BranchID:    branchID,
		TableID:     &tables[1].ID,
		OrderType:   "dine-in",
		TotalAmount: 35000,
		CreatedAt:   time.Now().Add(-time.Hour),
		Items: []domain.OrderItem{
			{MenuItemID: menuItemID, Quantity: 1, UnitPrice: 35000, Subtotal: 35000},
		},
	}
	if err := repo.Create(completedOrder); err != nil {
		log.Printf("seed completed order: %v", err)
	}

	if err := repo.UpdateStatus(completedOrder.ID, "completed"); err != nil {
		log.Printf("seed update order status: %v", err)
	}
	if err := repo.UpdatePayment(completedOrder.ID, "paid", "cash"); err != nil {
		log.Printf("seed update payment: %v", err)
	}

	log.Printf("created 2 sample orders")
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("seed hash password: %v", err)
		return password
	}
	return string(hash)
}
