package main

import (
	"log"
	"net/http"
	"os"

	"github.com/resto-fnb/backend/internal/config"
	"github.com/resto-fnb/backend/internal/handler/auth"
	"github.com/resto-fnb/backend/internal/handler/chain"
	"github.com/resto-fnb/backend/internal/handler/menu"
	"github.com/resto-fnb/backend/internal/handler/order"
	qrHandler "github.com/resto-fnb/backend/internal/handler/qr"
	"github.com/resto-fnb/backend/internal/handler/theme"
	"github.com/resto-fnb/backend/internal/handler"
	"github.com/resto-fnb/backend/internal/middleware"
	"github.com/resto-fnb/backend/internal/usecase"
	"github.com/resto-fnb/backend/internal/repository"
	authUsecase "github.com/resto-fnb/backend/internal/usecase/auth"
	chainUsecase "github.com/resto-fnb/backend/internal/usecase/chain"
	menuUsecase "github.com/resto-fnb/backend/internal/usecase/menu"
	orderUsecase "github.com/resto-fnb/backend/internal/usecase/order"
	qrUsecase "github.com/resto-fnb/backend/internal/usecase/qr"
	"github.com/resto-fnb/backend/pkg/database"
	"github.com/resto-fnb/backend/pkg/seed"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if os.Getenv("SEED_ON_START") == "true" {
		if err := seed.Run(db); err != nil {
			log.Printf("seed error: %v", err)
		}
	}

	chainRepo := repository.NewChainRepo(db)
	branchRepo := repository.NewBranchRepo(db)
	themeRepo := repository.NewThemeRepo(db)
	userRepo := repository.NewUserRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	menuItemRepo := repository.NewMenuItemRepo(db)
	variantRepo := repository.NewVariantRepo(db)
	availabilityRepo := repository.NewAvailabilityRepo(db)
	tableRepo := repository.NewTableRepo(db)
	orderRepo := repository.NewOrderRepo(db)
	carouselRepo := repository.NewCarouselRepo(db)

	chainUc := chainUsecase.NewUsecase(chainRepo)
	branchUc := chainUsecase.NewBranchUsecase(branchRepo, chainRepo)
	themeUc := chainUsecase.NewThemeUsecase(themeRepo, chainRepo)
	authUc := authUsecase.NewUsecase(userRepo, cfg.JWTSecret)
	categoryUc := menuUsecase.NewCategoryUsecase(categoryRepo)
	menuItemUc := menuUsecase.NewMenuItemUsecase(menuItemRepo)
	variantUc := menuUsecase.NewVariantUsecase(variantRepo)
	availabilityUc := menuUsecase.NewAvailabilityUsecase(availabilityRepo)
	tableUc := qrUsecase.NewTableUsecase(tableRepo, cfg.QRFrontendURL)
	orderUc := orderUsecase.NewUsecase(orderRepo)
	carouselUc := usecase.NewCarouselUsecase(carouselRepo)

	chainHandler := chain.NewHandler(chainUc)
	branchHandler := chain.NewBranchHandler(branchUc)
	themeHandler := theme.NewHandler(themeUc)
	authHandler := auth.NewHandler(authUc)
	categoryHandler := menu.NewCategoryHandler(categoryUc.Create, categoryUc.List, categoryUc.Update, categoryUc.Delete)
	menuItemHandler := menu.NewMenuItemHandler(menuItemUc.Create, menuItemUc.ListByChainID, menuItemUc.ListPublicByChainID, menuItemUc.Update, menuItemUc.Delete)
	variantHandler := menu.NewVariantHandler(variantUc.Create, variantUc.ListByMenuItemID, variantUc.Delete)
	availabilityHandler := menu.NewAvailabilityHandler(availabilityUc.SetAvailability, availabilityUc.GetByBranchID)
	tableHandler := qrHandler.NewHandler(tableUc, branchRepo)
	orderHandler := order.NewHandler(orderUc)
	carouselHandler := handler.NewCarouselHandler(
		carouselUc.Create, carouselUc.ListByChainID, carouselUc.ListPublic,
		carouselUc.Update, carouselUc.SoftDelete,
	)

	authMw := middleware.Auth(cfg.JWTSecret)
	adminMw := middleware.RequireRole("admin")
	cashierMw := middleware.RequireRole("admin", "cashier")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/refresh", authHandler.Refresh)

	mux.Handle("/api/admin/chains", authMw(adminMw(http.HandlerFunc(chainHandler.Register))))
	mux.Handle("/api/admin/branches", authMw(adminMw(http.HandlerFunc(branchHandler.Create))))
	mux.Handle("/api/admin/branches/list", authMw(adminMw(http.HandlerFunc(branchHandler.ListByChain))))
	mux.Handle("/api/admin/themes", authMw(adminMw(http.HandlerFunc(themeHandler.Upsert))))
	mux.Handle("/api/admin/themes/get", authMw(http.HandlerFunc(themeHandler.Get)))
	mux.Handle("/api/admin/staff/list", authMw(adminMw(http.HandlerFunc(authHandler.ListStaff))))
	mux.Handle("/api/admin/staff/deactivate", authMw(adminMw(http.HandlerFunc(authHandler.DeactivateStaff))))

	mux.Handle("/api/admin/categories/list", authMw(http.HandlerFunc(categoryHandler.List)))
	mux.Handle("/api/admin/categories/delete", authMw(adminMw(http.HandlerFunc(categoryHandler.Delete))))
	mux.Handle("/api/admin/categories/update", authMw(adminMw(http.HandlerFunc(categoryHandler.Update))))
	mux.Handle("/api/admin/categories", authMw(adminMw(http.HandlerFunc(categoryHandler.Create))))
	mux.Handle("/api/admin/menu-items/list", authMw(http.HandlerFunc(menuItemHandler.List)))
	mux.Handle("/api/admin/menu-items/delete", authMw(adminMw(http.HandlerFunc(menuItemHandler.Delete))))
	mux.Handle("/api/admin/menu-items/update", authMw(adminMw(http.HandlerFunc(menuItemHandler.Update))))
	mux.Handle("/api/admin/menu-items", authMw(adminMw(http.HandlerFunc(menuItemHandler.Create))))
	mux.Handle("/api/admin/variants", authMw(adminMw(http.HandlerFunc(variantHandler.Create))))
	mux.Handle("/api/admin/variants/list", authMw(http.HandlerFunc(variantHandler.List)))
	mux.Handle("/api/admin/availability", authMw(adminMw(http.HandlerFunc(availabilityHandler.Set))))
	mux.Handle("/api/admin/availability/list", authMw(http.HandlerFunc(availabilityHandler.Get)))

	mux.Handle("/api/admin/tables", authMw(adminMw(http.HandlerFunc(tableHandler.Create))))
	mux.Handle("/api/admin/tables/list", authMw(http.HandlerFunc(tableHandler.List)))
	mux.Handle("/api/admin/tables/update", authMw(adminMw(http.HandlerFunc(tableHandler.Update))))
	mux.Handle("/api/admin/tables/regenerate-qr", authMw(adminMw(http.HandlerFunc(tableHandler.RegenerateQr))))
	mux.Handle("/api/admin/tables/delete", authMw(adminMw(http.HandlerFunc(tableHandler.Delete))))

	mux.Handle("/api/upload", authMw(http.HandlerFunc(handler.UploadHandler)))

	mux.Handle("/api/admin/carousel", authMw(adminMw(http.HandlerFunc(carouselHandler.Create))))
	mux.Handle("/api/admin/carousel/list", authMw(http.HandlerFunc(carouselHandler.List)))
	mux.Handle("/api/admin/carousel/update", authMw(adminMw(http.HandlerFunc(carouselHandler.Update))))
	mux.Handle("/api/admin/carousel/delete", authMw(adminMw(http.HandlerFunc(carouselHandler.Delete))))
	mux.HandleFunc("/api/carousel/list", carouselHandler.ListPublic)

	mux.HandleFunc("/api/categories/list", categoryHandler.List)
	mux.HandleFunc("/api/menu-items/list", menuItemHandler.ListPublic)

	mux.HandleFunc("/api/orders", orderHandler.Create)
	mux.Handle("/api/orders/status", authMw(cashierMw(http.HandlerFunc(orderHandler.UpdateStatus))))
	mux.Handle("/api/orders/pay", authMw(cashierMw(http.HandlerFunc(orderHandler.Pay))))
	mux.Handle("/api/orders/by-branch", authMw(cashierMw(http.HandlerFunc(orderHandler.ListByBranch))))
	mux.Handle("/api/orders/by-chain", authMw(adminMw(http.HandlerFunc(orderHandler.ListByChain))))

	uploadPath, uploadHandler := handler.UploadRoutes()
	mux.Handle(uploadPath, uploadHandler)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
