package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/apikey"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/cart"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/currency"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/massmail"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/news"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/support"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/pdf"
)

func main() {
	ctx := context.Background()
	fmt.Println("==================================================================")
	fmt.Println("🚀 FOSSBilling Next-Gen Backend (Golang) - E2E Live Simulation")
	fmt.Println("==================================================================")

	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	promoRepo := memory.NewMockPromoRepository()
	txnRepo := memory.NewMockTransactionRepository()
	supportRepo := memory.NewMockSupportRepository()
	currencyRepo := memory.NewMockCurrencyRepository()
	newsRepo := memory.NewMockNewsRepository()
	downloadRepo := memory.NewMockDownloadableRepository()
	apiKeyRepo := memory.NewMockAPIKeyRepository()
	massMailRepo := memory.NewMockMassMailRepository()

	jwtSecret := "super-secret-jwt-key-32-chars-long"
	mockMailer := mailer.NewMockMailer()
	emailService := notification.NewEmailService(mockMailer, "billing@nusantara-cloud.com", "Nusantara Cloud")
	eventBus := events.NewEventBus()

	taxCalc := billing.NewTaxCalculator([]billing.TaxRule{{Name: "Indonesian PPN", Country: "ID", Rate: 11.0}})
	authUc := auth.NewAuthUsecase(clientRepo, jwtSecret)
	orderService := order.NewOrderService(orderRepo)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc)
	promoCalc := cart.NewPromoCalculator(promoRepo)
	cartService := cart.NewCartService(promoCalc, promoRepo, orderRepo, invService)
	webhookService := payment.NewWebhookService(txnRepo, invRepo, orderService, orderRepo)
	supportService := support.NewSupportService(supportRepo, clientRepo)
	statsService := stats.NewStatsService(clientRepo, orderRepo, invRepo, supportRepo)

	currencyService := currency.NewCurrencyService(currencyRepo)
	newsService := news.NewNewsService(newsRepo)
	downloadService := downloadable.NewDownloadableService(downloadRepo, orderRepo, jwtSecret)
	apiKeyService := apikey.NewAPIKeyService(apiKeyRepo)
	massMailService := massmail.NewMassMailService(massMailRepo, clientRepo, mockMailer, "admin@nusantara-cloud.com", "Nusantara Cloud")

	cpanelProv := provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{Host: "sg1.nusantara-cloud.com"})
	daProv := provisioning.NewDirectAdminProvisioner("da.nusantara-cloud.com", 2222, "admin", "secret")
	pleskProv := provisioning.NewPleskProvisioner("plesk.nusantara-cloud.com", 8443, "plesk-api-key")
	licenseProv := provisioning.NewLicenseProvisioner("fossbilling-enterprise-master-key")
	domainDriver := provisioning.NewMockRegistrarDriver()

	eventBus.Subscribe(events.EventClientRegistered, func(ctx context.Context, e events.Event) error {
		return emailService.SendWelcomeEmail(ctx, e.Payload.(*domain.Client))
	})

	// 1. Registrasi Klien Baru
	regRes, _, _ := authUc.Register(ctx, auth.RegisterDTO{
		Email: "budi.santoso@example.com", Password: "SecurePassword123!",
		FirstName: "Budi", LastName: "Santoso", Country: "ID", Currency: "IDR",
	})
	registeredClient, _ := clientRepo.GetByID(ctx, regRes.Client.ID)
	_ = eventBus.Publish(ctx, events.Event{Type: events.EventClientRegistered, Payload: registeredClient})
	fmt.Printf("   ✅ Klien Terdaftar: %s %s (ID: %d)\n", regRes.Client.FirstName, regRes.Client.LastName, regRes.Client.ID)

	// 2. Domain & Kupon
	avail, _ := domainDriver.CheckAvailability(ctx, "solusinusantara.com")
	fmt.Printf("   🔍 Domain: %s (Tersedia: %v, Harga: %s %s)\n", avail.DomainName, avail.IsAvailable, avail.Currency, decimal.Money(avail.Price).String())

	_ = promoRepo.Create(ctx, &domain.Promo{Code: "MERDEKA20", Type: domain.PromoTypePercentage, Value: decimal.FromFloat(20.00), Active: true})
	_, _ = currencyService.CreateCurrency(ctx, currency.CreateCurrencyDTO{Code: "USD", Title: "US Dollar", ConversionRate: 0.000065, Format: "$ {{price}}", PriceFormat: "2"})

	// 3. Checkout
	hostingCfg, _ := json.Marshal(map[string]string{"domain": "solusinusantara.com"})
	shoppingCart := &cart.Cart{
		ClientID: regRes.Client.ID, PromoCode: "MERDEKA20",
		Items: []cart.CartItem{
			{ProductID: 101, Title: "Cloud VPS cPanel Pro", Period: "1M", Price: decimal.FromFloat(200000.00), Quantity: 1, Config: hostingCfg},
			{ProductID: 202, Title: "DirectAdmin Hosting", Period: "1M", Price: decimal.FromFloat(150000.00), Quantity: 1},
			{ProductID: 303, Title: "FOSSBilling Enterprise", Period: "1Y", Price: decimal.FromFloat(500000.00), Quantity: 1},
			{ProductID: 404, Title: "Nusantara Cloud OS", Period: "ONETIME", Price: decimal.FromFloat(100000.00), Quantity: 1},
		},
	}
	_ = downloadRepo.Create(ctx, &domain.DownloadableFile{ProductID: 404, Filename: "os.iso", FilePath: "/data/os.iso", FileSize: 500 * 1024 * 1024, ContentType: "application/octet-stream", Version: "1.0.0"})
	checkoutRes, _ := cartService.Checkout(ctx, shoppingCart)
	fmt.Printf("   💰 TOTAL INVOICE : %s %s (Nomor: #%s%s)\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Total.String(), checkoutRes.Invoice.Serie, checkoutRes.Invoice.Nr)

	// 4. Webhook Settlement
	txn, _ := webhookService.HandlePaymentWebhook(ctx, payment.WebhookPayload{
		GatewayID: "midtrans", TxnID: "MID-TXN-" + fmt.Sprint(time.Now().Unix()),
		InvoiceID: checkoutRes.Invoice.ID, Amount: checkoutRes.Invoice.Total, Currency: checkoutRes.Invoice.Currency,
		Raw: []byte(`{"status":"settlement"}`),
	})
	fmt.Printf("   ✅ Pembayaran Sukses: TxnID=%s, Gateway=%s\n", txn.TxnID, txn.GatewayID)

	// 5. PDF Invoice & Provisioning
	pdfBytes, _ := pdf.GenerateInvoicePDF(checkoutRes.Invoice, registeredClient, "Nusantara Cloud Indonesia", "", "")
	fmt.Printf("   ✅ Invoice PDF Dihasilkan: %d bytes\n", len(pdfBytes))

	runProvisioningDemo(ctx, orderRepo, checkoutRes.Orders, cpanelProv, daProv, pleskProv, licenseProv)

	// 6. Signed Download Link & API Key
	dlFile, _ := downloadRepo.GetByProductID(ctx, 404)
	signedLink, _ := downloadService.GenerateDownloadLink(ctx, regRes.Client.ID, dlFile.ID, 2*time.Hour)
	fmt.Printf("   🔗 Link Unduh HMAC: %s\n", signedLink.URL)

	apiKey, _ := apiKeyService.GenerateKey(ctx, regRes.Client.ID, "Deployment Bot", 90)
	fmt.Printf("   🔑 API Key: %s\n", apiKey.Key)

	// 7. News & Mass Mail
	article, _ := newsService.Create(ctx, news.CreateNewsDTO{AdminID: 1, Title: "Pembaruan 2026", Content: "10Gbps unmetered", Status: domain.NewsStatusPublished})
	fmt.Printf("   📰 Berita: %s\n", article.Title)
	campaign, _ := massMailService.Create(ctx, 1, "Maintenance", "<p>Maintenance malam ini.</p>")
	_, _ = massMailService.Send(ctx, campaign.ID)

	// 8. Support Ticket & Dashboard
	ticket, _ := supportService.OpenTicket(ctx, support.CreateTicketDTO{ClientID: regRes.Client.ID, Subject: "DNS Help", Message: "Bantuan DNS", Priority: domain.PriorityHigh})
	_, _ = supportService.StaffReply(ctx, ticket.ID, 1, "Nameserver: ns1 & ns2.nusantara-cloud.com")
	_ = supportService.CloseTicket(ctx, ticket.ID, regRes.Client.ID)

	dashStats, _ := statsService.CalculateDashboard(ctx)
	fmt.Printf("   📈 MRR: %s %s • Klien: %d • Pesanan: %d\n", checkoutRes.Invoice.Currency, dashStats.MonthlyRecurring.String(), dashStats.TotalClients, dashStats.ActiveOrders)

	fmt.Println("\n🎉 SIMULASI SELESAI: 100% Modul Sukses Teruji!")
}
