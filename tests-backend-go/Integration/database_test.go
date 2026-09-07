package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getDBPool(t *testing.T) *pgxpool.Pool {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://fossbilling:secretpassword123@localhost:5433/fossbilling?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping Database Integration Test: Cannot connect to PostgreSQL at %s (%v)", dbURL, err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("Skipping Database Integration Test: PostgreSQL not reachable at %s (%v)", dbURL, err)
		return nil
	}

	return pool
}

// 1. TDD: Test that all 19 required database tables exist in the schema
func TestDatabaseSchema_AllTablesExist(t *testing.T) {
	pool := getDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	ctx := context.Background()

	expectedTables := []string{
		"admin_groups",
		"staff",
		"clients",
		"client_balances",
		"products",
		"client_orders",
		"invoices",
		"invoice_items",
		"transactions",
		"support_tickets",
		"support_ticket_messages",
		"currencies",
		"promos",
		"promo_redemptions",
		"news_posts",
		"downloadable_files",
		"api_keys",
		"audit_logs",
		"activity_logs",
		"mass_mail_campaigns",
		"redirects",
		"extensions",
		"forms",
		"form_fields",
		"blocked_ips",
		"servers",
	}

	for _, tbl := range expectedTables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`
		if err := pool.QueryRow(ctx, query, tbl).Scan(&exists); err != nil {
			t.Fatalf("Failed to check table %s: %v", tbl, err)
		}
		if !exists {
			t.Errorf("Expected table %q to exist in database, but it was missing", tbl)
		}
	}
}

// 2. TDD: Test Table Column Constraints & Foreign Keys
func TestDatabaseSchema_ForeignKeyConstraints(t *testing.T) {
	pool := getDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	ctx := context.Background()

	// Check that client_orders.product_id is NULLABLE
	var isNullable string
	query := `SELECT is_nullable 
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		AND table_name = 'client_orders' 
		AND column_name = 'product_id'`
	if err := pool.QueryRow(ctx, query).Scan(&isNullable); err != nil {
		t.Fatalf("Failed to check is_nullable for client_orders.product_id: %v", err)
	}

	if isNullable != "YES" {
		t.Errorf("client_orders.product_id must be nullable (YES) to support custom/domain items, got %s", isNullable)
	}

	// Check that invoice_items cascade deletes on invoice removal
	var count int
	fkQuery := `SELECT count(*) 
		FROM information_schema.referential_constraints rc
		JOIN information_schema.table_constraints tc ON rc.constraint_name = tc.constraint_name
		WHERE tc.table_name = 'invoice_items' AND rc.delete_rule = 'CASCADE'`
	if err := pool.QueryRow(ctx, fkQuery).Scan(&count); err != nil {
		t.Fatalf("Failed to check cascade FK on invoice_items: %v", err)
	}
	if count == 0 {
		t.Errorf("Expected cascade delete constraint on invoice_items referencing invoices")
	}
}

// 3. TDD: Test Seed Data Verification
func TestDatabaseSchema_SeedDataVerification(t *testing.T) {
	pool := getDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	ctx := context.Background()

	// Verify Admin Staff
	var adminEmail string
	err := pool.QueryRow(ctx, `SELECT email FROM staff WHERE email = 'admin@fossbilling.org'`).Scan(&adminEmail)
	if err != nil {
		t.Errorf("Default superadmin 'admin@fossbilling.org' missing: %v", err)
	}

	// Verify Default Currency USD
	var usdCount int
	err = pool.QueryRow(ctx, `SELECT count(*) FROM currencies WHERE code = 'USD' AND is_default = TRUE`).Scan(&usdCount)
	if err != nil || usdCount != 1 {
		t.Errorf("Default currency USD missing or not marked as default: %v", err)
	}

	// Verify Default Products
	var productCount int
	err = pool.QueryRow(ctx, `SELECT count(*) FROM products WHERE status = 'enabled'`).Scan(&productCount)
	if err != nil || productCount < 5 {
		t.Errorf("Expected at least 5 default seeded products, got %d (err: %v)", productCount, err)
	}

	// Verify Default Promo
	var promoCode string
	err = pool.QueryRow(ctx, `SELECT code FROM promos WHERE code = 'MERDEKA20' AND active = TRUE`).Scan(&promoCode)
	if err != nil {
		t.Errorf("Default promo 'MERDEKA20' missing: %v", err)
	}
}

// 4. TDD: Test ClientRepository with Balance Ledger
func TestDatabase_ClientRepositoryCRUD(t *testing.T) {
	pool := getDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)

	uniqueEmail := fmt.Sprintf("test.client.%d@example.com", time.Now().UnixNano())
	client := &domain.Client{
		Email:        uniqueEmail,
		PasswordHash: "hashed-pw-test",
		FirstName:    "John",
		LastName:     "Doe",
		Country:      "US",
		Currency:     "USD",
		Status:       "active",
	}

	// Create
	if err := clientRepo.Create(ctx, client); err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if client.ID <= 0 {
		t.Fatalf("Expected valid client ID, got %d", client.ID)
	}

	// Read by ID
	fetched, err := clientRepo.GetByID(ctx, client.ID)
	if err != nil {
		t.Fatalf("Failed to get client by ID: %v", err)
	}
	if fetched.Email != uniqueEmail || fetched.FirstName != "John" {
		t.Errorf("Fetched client data mismatch: %+v", fetched)
	}

	// Read by Email
	fetchedByEmail, err := clientRepo.GetByEmail(ctx, uniqueEmail)
	if err != nil || fetchedByEmail.ID != client.ID {
		t.Fatalf("Failed to get client by Email: %v", err)
	}

	// Adjust Balance (Credit $50.00)
	creditAmt := decimal.FromFloat(50.00)
	if err := clientRepo.AddBalanceTransaction(ctx, &domain.ClientBalance{
		ClientID:    client.ID,
		Type:        "credit",
		Amount:      creditAmt,
		Description: "Welcome Deposit",
	}); err != nil {
		t.Fatalf("Failed to add balance transaction: %v", err)
	}

	balance, err := clientRepo.GetBalance(ctx, client.ID)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}
	if balance != creditAmt {
		t.Errorf("Expected balance %v, got %v", creditAmt, balance)
	}
}

// 5. TDD: Test OrderRepository with Valid Product and Null Product
func TestDatabase_OrderRepositoryWithNullProduct(t *testing.T) {
	pool := getDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	orderRepo := postgres.NewOrderRepository(pool)

	// Create client
	client := &domain.Client{
		Email:        fmt.Sprintf("order.test.%d@example.com", time.Now().UnixNano()),
		PasswordHash: "secret",
		FirstName:    "Order",
		LastName:     "Tester",
		Country:      "ID",
		Currency:     "USD",
		Status:       "active",
	}
	if err := clientRepo.Create(ctx, client); err != nil {
		t.Fatalf("Client creation failed: %v", err)
	}

	// 1. Create order with existing product_id (e.g. 1)
	validOrder := &domain.Order{
		ClientID:  client.ID,
		ProductID: 1,
		Title:     "cPanel Starter Cloud",
		Period:    "1M",
		Price:     decimal.FromFloat(9.99),
		Currency:  "USD",
		Status:    domain.OrderStatusPendingSetup,
	}
	if err := orderRepo.Create(ctx, validOrder); err != nil {
		t.Fatalf("Failed to create order with valid product: %v", err)
	}
	if validOrder.ID <= 0 {
		t.Fatalf("Expected valid order ID, got %d", validOrder.ID)
	}

	// 2. Create order with NO product (domain registration / custom item, ProductID: 0)
	domainOrder := &domain.Order{
		ClientID:  client.ID,
		ProductID: 0,
		Title:     "Domain Registration: mycompany.com",
		Period:    "1Y",
		Price:     decimal.FromFloat(12.00),
		Currency:  "USD",
		Status:    domain.OrderStatusPendingSetup,
	}
	if err := orderRepo.Create(ctx, domainOrder); err != nil {
		t.Fatalf("Failed to create domain order without product_id: %v", err)
	}
	if domainOrder.ID <= 0 {
		t.Fatalf("Expected valid domain order ID, got %d", domainOrder.ID)
	}

	// 3. Fetch orders for client
	orders, total, err := orderRepo.ListByClientID(ctx, client.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list orders: %v", err)
	}
	if total != 2 || len(orders) != 2 {
		t.Fatalf("Expected 2 orders, got total=%d len=%d", total, len(orders))
	}

	// 4. Update status transition
	if err := orderRepo.UpdateStatus(ctx, validOrder.ID, domain.OrderStatusActive, nil); err != nil {
		t.Fatalf("Failed to update status to active: %v", err)
	}
	updated, err := orderRepo.GetByID(ctx, validOrder.ID)
	if err != nil || updated.Status != domain.OrderStatusActive {
		t.Errorf("Expected status active, got %s", updated.Status)
	}
}

// 6. TDD: Test Full Cart Checkout Flow with Real Database Persistence
func TestDatabase_E2ECheckoutFlow(t *testing.T) {
	pool := getDBPool(t)
	if pool == nil {
		return
	}
	defer pool.Close()

	ctx := context.Background()
	clientRepo := postgres.NewClientRepository(pool)
	orderRepo := postgres.NewOrderRepository(pool)
	invRepo := postgres.NewInvoiceRepository(pool)
	promoRepo := postgres.NewPromoRepository(pool)

	taxCalc := billing.NewTaxCalculator([]billing.TaxRule{{Name: "VAT", Country: "US", Rate: 0.0}})
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc, nil)
	promoCalc := cart.NewPromoCalculator(promoRepo)
	cartService := cart.NewCartService(promoCalc, promoRepo, orderRepo, clientRepo, taxCalc, invService)

	// Create test client
	client := &domain.Client{
		Email:        fmt.Sprintf("checkout.db.%d@example.com", time.Now().UnixNano()),
		PasswordHash: "pass123",
		FirstName:    "Budi",
		LastName:     "Pratama",
		Country:      "US",
		Currency:     "USD",
		Status:       "active",
	}
	if err := clientRepo.Create(ctx, client); err != nil {
		t.Fatalf("Client create failed: %v", err)
	}

	shoppingCart := &cart.Cart{
		ClientID:  client.ID,
		PromoCode: "WELCOME10",
		Items: []cart.CartItem{
			{
				ProductID: 1,
				Title:     "cPanel Starter Cloud",
				Period:    "1M",
				Price:     decimal.FromFloat(10.00),
				Quantity:  1,
			},
			{
				ProductID: 99, // Domain registration
				Title:     "Domain Registration: myawesomebrand.com",
				Period:    "1Y",
				Price:     decimal.FromFloat(15.00),
				Quantity:  1,
			},
		},
	}

	res, err := cartService.Checkout(ctx, shoppingCart)
	if err != nil {
		t.Fatalf("Cart checkout failed: %v", err)
	}

	if len(res.Orders) != 2 {
		t.Errorf("Expected 2 orders created, got %d", len(res.Orders))
	}
	if res.Invoice == nil || res.Invoice.ID <= 0 {
		t.Fatalf("Expected valid generated invoice, got %+v", res.Invoice)
	}

	// Verify Invoice was stored in DB with items
	dbInvoice, err := invRepo.GetByID(ctx, res.Invoice.ID)
	if err != nil {
		t.Fatalf("Failed to fetch generated invoice from DB: %v", err)
	}
	if len(dbInvoice.Items) < 2 {
		t.Errorf("Expected at least 2 items on invoice, got %d", len(dbInvoice.Items))
	}
}
