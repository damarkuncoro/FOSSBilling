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
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getDBPool(t *testing.T) *pgxpool.Pool {
	u := os.Getenv("DATABASE_URL"); if u == "" { u = "postgres://fossbilling:secretpassword123@localhost:5433/fossbilling?sslmode=disable" }
	ctx, c := context.WithTimeout(context.Background(), 5*time.Second); defer c()
	p, err := pgxpool.New(ctx, u); if err != nil { t.Skip("No DB"); return nil }; if p.Ping(ctx) != nil { p.Close(); t.Skip("DB Down"); return nil }; return p
}

func TestDatabase_ClientRepositoryCRUD(t *testing.T) {
	p := getDBPool(t); if p == nil { return }; defer p.Close(); ctx := context.Background(); r := postgres.NewClientRepository(p)
	em := fmt.Sprintf("t.%d@e.com", time.Now().UnixNano()); c := &domain.Client{Email: em, PasswordHash: "h", FirstName: "J", LastName: "D", Country: "US", Currency: "USD", Status: "active"}
	if err := r.Create(ctx, c); err != nil { t.Fatalf("Create failed: %v", err) }
	cr := decimal.FromFloat(50.00); _ = r.AddBalanceTransaction(ctx, &domain.ClientBalance{ClientID: c.ID, Type: "credit", Amount: cr, Description: "W"})
	bal, _ := r.GetBalance(ctx, c.ID); if bal != cr { t.Errorf("Balance mismatch: %v != %v", bal, cr) }
}

func TestDatabase_E2ECheckoutFlow(t *testing.T) {
	p := getDBPool(t); if p == nil { return }; defer p.Close(); ctx := context.Background()
	cr, or, ir, pmr, pr, cor := postgres.NewClientRepository(p), postgres.NewOrderRepository(p), postgres.NewInvoiceRepository(p), postgres.NewPromoRepository(p), postgres.NewProductRepository(p), postgres.NewCompanyRepository(p)
	tc := billing.NewTaxCalculator(nil); is := billing.NewInvoiceService(ir, cr, cor, tc, nil); pc := cart.NewPromoCalculator(pmr); fs := formbuilder.NewFormbuilderService(postgres.NewFormbuilderRepository(p))
	cs := cart.NewCartService(pc, pmr, or, pr, cr, fs, tc, is, nil, nil)
	em := fmt.Sprintf("c.%d@e.com", time.Now().UnixNano()); c := &domain.Client{Email: em, PasswordHash: "p", FirstName: "B", LastName: "P", Country: "US", Currency: "USD", Status: "active"}
	_ = cr.Create(ctx, c); _ = pr.Create(ctx, &domain.Product{ID: 1, Title: "P1", PriceMonthly: 100000})
	sc := &cart.Cart{ClientID: c.ID, Items: []cart.CartItem{{ProductID: 1, Title: "P1", Period: "1M", Price: 100000, Quantity: 1}}}
	res, err := cs.Checkout(ctx, sc, "1.1.1.1"); if err != nil { t.Fatalf("Checkout failed: %v", err) }
	if len(res.Orders) != 1 || res.Invoice == nil { t.Error("Failed to create orders/invoice") }
}
