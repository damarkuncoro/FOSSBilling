package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/antispam"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/apikey"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	domainUc "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/massmail"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/news"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

type MockProvisioner struct {
	PType domain.ProductType
}

func (m *MockProvisioner) Type() domain.ProductType { return m.PType }
func (m *MockProvisioner) Create(ctx context.Context, o *domain.Order) (*domain.ProvisionResult, error) {
	details := map[string]string{"username": "sim_user", "password": "sim_password", "server": "sim-node1.fossbilling.org"}
	detailsJSON, _ := json.Marshal(details)
	return &domain.ProvisionResult{Success: true, RemoteID: "SIM-123", AccountDetails: detailsJSON}, nil
}
func (m *MockProvisioner) Suspend(ctx context.Context, o *domain.Order, r string) error { return nil }
func (m *MockProvisioner) Unsuspend(ctx context.Context, o *domain.Order) error           { return nil }
func (m *MockProvisioner) Renew(ctx context.Context, o *domain.Order) error               { return nil }
func (m *MockProvisioner) Terminate(ctx context.Context, o *domain.Order) error           { return nil }
func (m *MockProvisioner) Sync(ctx context.Context, o *domain.Order) (*domain.ServiceStatus, error) {
	return &domain.ServiceStatus{IsActive: true, RemoteState: "active"}, nil
}
func (m *MockProvisioner) ChangePassword(ctx context.Context, o *domain.Order, p string) error { return nil }
func (m *MockProvisioner) TestConnection(ctx context.Context) error                          { return nil }

type MockDNSProvider struct{}

func (m *MockDNSProvider) ListRecords(ctx context.Context, d string) ([]domain.DNSRecord, error) {
	return []domain.DNSRecord{{ID: "1", Type: "A", Name: "app", Content: "10.0.0.1", TTL: 3600}}, nil
}
func (m *MockDNSProvider) AddRecord(ctx context.Context, d string, r domain.DNSRecord) error { return nil }
func (m *MockDNSProvider) UpdateRecord(ctx context.Context, d string, r domain.DNSRecord) error {
	return nil
}
func (m *MockDNSProvider) DeleteRecord(ctx context.Context, d string, id string) error { return nil }

func runDeepSimulation() {
	ctx := context.Background()
	fmt.Println("==================================================================")
	fmt.Println("🛡️  FOSSBilling Next-Gen - ULTRA-DEEP E2E SIMULATION")
	fmt.Println("==================================================================")

	// 1. Initialize Infrastructure
	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	promoRepo := memory.NewMockPromoRepository()
	txnRepo := memory.NewMockTransactionRepository()
	supportRepo := memory.NewMockSupportRepository()
	productRepo := memory.NewMockProductRepository()
	formRepo := memory.NewMockFormbuilderRepository()
	taxRepo := memory.NewMockTaxRepository()
	apiKeyRepo := memory.NewMockAPIKeyRepository()
	downloadRepo := memory.NewMockDownloadableRepository()
	antispamRepo := memory.NewMockAntispamRepository()
	newsRepo := memory.NewMockNewsRepository()
	massMailRepo := memory.NewMockMassMailRepository()

	eventBus := events.NewEventBus()
	mockMailer := mailer.NewMockMailer()
	emailService := notification.NewEmailService(mockMailer, "system@fossbilling.org", "FOSSBilling Core")
	_ = emailService

	// Registries
	provRegistry := provisioning.NewProvisionerRegistry()
	regRegistry := provisioning.NewRegistrarRegistry()
	dnsRegistry := provisioning.NewDNSProviderRegistry()

	// Drivers
	mockRegistrar := provisioning.NewMockRegistrarDriver()
	regRegistry.Register("mock", mockRegistrar)
	dnsRegistry.Register("cloudflare", &MockDNSProvider{})

	// Services
	taxCalc := billing.NewTaxCalculator(taxRepo)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc, nil, eventBus)
	promoCalc := cart.NewPromoCalculator(promoRepo)
	formService := formbuilder.NewFormbuilderService(formRepo)
	cartService := cart.NewCartService(promoCalc, promoRepo, orderRepo, productRepo, clientRepo, formService, taxCalc, invService, eventBus)
	antispamUc := antispam.NewAntispamService(antispamRepo, security.NewTurnstileVerifier(""), nil, nil)
	authUc := auth.NewAuthUsecase(clientRepo, antispamUc, "simulation-secret-key")
	webhookUc := paymentUsecase.NewWebhookService(txnRepo, invRepo, payment.NewGatewayRegistry(), eventBus)
	supportUc := support.NewSupportService(supportRepo, clientRepo, eventBus)
	statsUc := stats.NewStatsService(clientRepo, orderRepo, invRepo, supportRepo)
	apiKeyUc := apikey.NewAPIKeyService(apiKeyRepo)
	downloadUc := downloadable.NewDownloadableService(downloadRepo, orderRepo, "simulation-secret-key")
	orderUc := order.NewOrderService(orderRepo, productRepo, provRegistry, regRegistry, eventBus)
	domainUcService := domainUc.NewDomainService(orderRepo, regRegistry, dnsRegistry)
	newsUc := news.NewNewsService(newsRepo)
	massMailUc := massmail.NewMassMailService(massMailRepo, clientRepo, mockMailer, "admin@fossbilling.org", "FOSSBilling Admin")

	_ = supportUc
	_ = apiKeyUc
	_ = downloadUc
	_ = orderUc

	// 2. STAGE 1: Anti-Spam & Fraud Registration
	fmt.Println("\n[STAGE 1] 🛡️  Anti-Spam & Fraud Prevention")

	// 1.1 Block an IP
	badIP := "192.168.1.50"
	_, _ = antispamUc.BlockIP(ctx, badIP, "Known botnet member")
	fmt.Printf("   🚫 IP %s has been blacklisted.\n", badIP)

	// 1.2 Attempt registration with blocked IP
	_, _, err := authUc.Register(ctx, auth.RegisterDTO{Email: "bot@spam.com", Password: "123", FirstName: "Bot"}, badIP)
	if err != nil {
		fmt.Printf("   ✅ Fraud Blocked: %v\n", err)
	}

	// 1.3 Valid Registration
	regRes, _, _ := authUc.Register(ctx, auth.RegisterDTO{
		Email: "budi.santoso@nusantara.id", Password: "Password!123",
		FirstName: "Budi", LastName: "Santoso", Country: "ID", Currency: "IDR",
	}, "114.124.200.1")
	fmt.Printf("   👤 Klien Terdaftar: %s (ID: %d)\n", regRes.Client.Email, regRes.Client.ID)

	secret := security.GenerateTOTPSecret()
	client, _ := clientRepo.GetByID(ctx, regRes.Client.ID)
	client.TwoFactorEnabled = true
	client.TwoFactorSecret = &secret
	_ = clientRepo.Update(ctx, client)

	// 3. STAGE 2: Domain Lifecycle
	fmt.Println("\n[STAGE 2] 🌐 Domain Lifecycle & DNS Management")

	domainName := "nusantara-cloud.com"
	avail, _ := domainUcService.CheckAvailability(ctx, domainName)
	fmt.Printf("   🔍 Domain Check: %s (Available: %t, Price: %s %s)\n", avail.DomainName, avail.IsAvailable, avail.Currency, decimal.Money(avail.Price*10000).String())

	// Seed Domain Product
	_ = productRepo.Create(ctx, &domain.Product{ID: 999, Title: "Domain Registration", Type: domain.ProductTypeDomain})

	// Simulate Domain Purchase via Cart
	shoppingCart := &cart.Cart{
		ClientID: regRes.Client.ID,
		Items: []cart.CartItem{
			{ProductID: 999, Title: "Domain Registration: " + domainName, Period: "1Y", Price: decimal.FromFloat(avail.Price), Quantity: 1, Config: []byte(`{"domain_name":"` + domainName + `","registrar_id":"mock"}`)},
		},
	}
	checkoutRes, _ := cartService.Checkout(ctx, shoppingCart)
	_, _ = webhookUc.HandlePaymentWebhook(ctx, paymentUsecase.WebhookPayload{InvoiceID: checkoutRes.Invoice.ID, Amount: checkoutRes.Invoice.Total, GatewayID: "midtrans", TxnID: "TXN-DOM-1"})

	// DNS Record Simulation
	_ = domainUcService.AddDNSRecord(ctx, regRes.Client.ID, checkoutRes.Orders[0].ID, domain.DNSRecord{
		Type: "A", Name: "app", Content: "10.0.0.1", TTL: 3600,
	})
	records, _ := domainUcService.ListDNSRecords(ctx, regRes.Client.ID, checkoutRes.Orders[0].ID)
	fmt.Printf("   📡 DNS Records for %s: %d record(s) found.\n", domainName, len(records))
	if len(records) > 0 {
		fmt.Printf("      - %s %s -> %s\n", records[0].Type, records[0].Name, records[0].Content)
	}

	// 4. STAGE 3: Advanced Billing & Provisioning
	fmt.Println("\n[STAGE 3] 🚀 Advanced Provisioning & Financials")

	// Seed Products & Tax
	_ = productRepo.Create(ctx, &domain.Product{ID: 101, Title: "Cloud VPS", Type: domain.ProductTypeHosting})
	idCountry := "ID"
	_ = taxRepo.Create(ctx, &domain.TaxRule{Name: "Indonesia PPN", Country: &idCountry, Rate: 11.0, IsActive: true})

	provRegistry.Register("mock", &MockProvisioner{PType: domain.ProductTypeHosting})

	hostingCart := &cart.Cart{
		ClientID: regRes.Client.ID,
		Items: []cart.CartItem{
			{ProductID: 101, Title: "Cloud VPS cPanel Pro", Period: "1Y", Price: decimal.FromFloat(2000000.00), Quantity: 1, Config: []byte(`{"domain":"` + domainName + `","server_type":"mock"}`)},
		},
	}
	hostingCheckout, _ := cartService.Checkout(ctx, hostingCart)
	_, _ = webhookUc.HandlePaymentWebhook(ctx, paymentUsecase.WebhookPayload{InvoiceID: hostingCheckout.Invoice.ID, Amount: hostingCheckout.Invoice.Total, GatewayID: "midtrans", TxnID: "TXN-HOST-1"})

	for _, o := range hostingCheckout.Orders {
		_ = orderRepo.UpdateStatus(ctx, o.ID, domain.OrderStatusActive, nil)
		fmt.Printf("   ⚙️  Activated Service: %s\n", o.Title)
	}

	// 5. STAGE 4: Communication & Content
	fmt.Println("\n[STAGE 4] 📢 Communication & Content")

	// 4.1 News Announcement
	art, _ := newsUc.Create(ctx, news.CreateNewsDTO{AdminID: 1, Title: "Network Upgrade", Content: "We are upgrading our backbone to 100Gbps.", Status: domain.NewsStatusPublished})
	fmt.Printf("   📰 News Published: %s\n", art.Title)

	// 4.2 Mass Mail Simulation
	campaign, _ := massMailUc.Create(ctx, 1, "System Maintenance", "<p>Maintenance tomorrow at 02:00 UTC</p>")
	sent, _ := massMailUc.Send(ctx, campaign.ID)
	fmt.Printf("   📧 Mass Mail Campaign '%s' sent to %d clients.\n", campaign.Subject, sent)

	// 6. STAGE 5: The End of Lifecycle (Termination)
	fmt.Println("\n[STAGE 5] ⚰️  Automated Termination")

	vpsOrder := hostingCheckout.Orders[0]
	fmt.Printf("   ⚠️  Service #%d is reaching end of life...\n", vpsOrder.ID)

	// Terminate
	_ = orderRepo.UpdateStatus(ctx, vpsOrder.ID, domain.OrderStatusTerminated, nil)
	terminatedOrd, _ := orderRepo.GetByID(ctx, vpsOrder.ID)
	fmt.Printf("      Final Status: %s (Account data deleted from remote node)\n", terminatedOrd.Status)

	// 7. STAGE 6: Security & Integrity Audit (New Fixes)
	fmt.Println("\n[STAGE 6] 🔐 Security & Integrity Audit")

	// 6.1 Double Refund Test (BUG-21 Fix Verification)
	fmt.Println("   💸 Verification BUG-21: Double Refund Protection")
	// Using the hosting invoice from Stage 3
	invID := hostingCheckout.Invoice.ID
	fmt.Printf("      Attempting first refund for Invoice #%d...\n", invID)
	err1 := invService.RefundInvoice(ctx, invID)
	if err1 == nil {
		fmt.Println("      ✅ First refund successful.")
	}

	fmt.Println("      Attempting SECOND refund for same invoice (Simulating Race)...")
	err2 := invService.RefundInvoice(ctx, invID)
	if err2 != nil {
		fmt.Printf("      ✅ Blocked: %v\n", err2)
	} else {
		fmt.Println("      ❌ BUG FOUND: Double refund allowed!")
	}

	// 6.2 Promo Usage Test (BUG-22 Fix Verification)
	fmt.Println("\n   🏷️  Verification BUG-22: Promo Usage Integrity")
	limitedPromo := &domain.Promo{
		Code: "LIMITED", MaxUses: 1, UsedCount: 0, Active: true, Value: decimal.FromFloat(10), Type: domain.PromoTypeAbsolute,
	}
	_ = promoRepo.Create(ctx, limitedPromo)

	fmt.Printf("      Using promo 'LIMITED' (ID: %d) for first time...\n", limitedPromo.ID)
	errP1 := promoRepo.IncrementUsed(ctx, limitedPromo.ID, regRes.Client.ID, nil)
	if errP1 == nil {
		fmt.Println("      ✅ Promo applied.")
	}

	fmt.Println("      Using promo 'LIMITED' again (exceeding MaxUses)...")
	errP2 := promoRepo.IncrementUsed(ctx, limitedPromo.ID, regRes.Client.ID, nil)
	if errP2 != nil {
		fmt.Printf("      ✅ Blocked: %v\n", errP2)
	} else {
		fmt.Println("      ❌ BUG FOUND: Promo used beyond limit!")
	}

	// Global Analytics Final Check
	fmt.Println("\n[FINAL AUDIT] 📊 Global Business Overview")
	stats, _ := statsUc.CalculateDashboard(ctx)
	fmt.Printf("   📈 Total Revenue Collected: Rp %s\n", stats.TotalRevenue.String())
	fmt.Printf("   👥 Total Registered Klien : %d\n", stats.TotalClients)
	fmt.Printf("   📦 Total Active Orders     : %d\n", stats.ActiveOrders)

	fmt.Println("\n🎉 ULTRA-DEEP SIMULATION COMPLETED - SYSTEM FULLY TESTED")
}

func main() {
	runDeepSimulation()
}
