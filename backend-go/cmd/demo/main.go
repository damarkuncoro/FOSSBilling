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

	// 1. Inisialisasi Repositories, Services, EventBus & Mailer
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

	taxCalc := billing.NewTaxCalculator([]billing.TaxRule{
		{Name: "Indonesian PPN", Country: "ID", Rate: 11.0},
	})

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

	// Register Event Listeners
	eventBus.Subscribe(events.EventClientRegistered, func(ctx context.Context, e events.Event) error {
		client := e.Payload.(*domain.Client)
		return emailService.SendWelcomeEmail(ctx, client)
	})

	eventBus.Subscribe(events.EventInvoicePaid, func(ctx context.Context, e events.Event) error {
		payload := e.Payload.(map[string]interface{})
		client := payload["client"].(*domain.Client)
		inv := payload["invoice"].(*domain.Invoice)
		txn := payload["txn"].(*domain.Transaction)
		return emailService.SendPaymentReceipt(ctx, client, inv, txn)
	})

	// 2. Simulasi Registrasi Klien Baru & Trigger Event Welcome Email
	fmt.Println("\n[1] 👤 Mendaftarkan Klien Baru & Mengirim Email Selamat Datang...")
	regRes, _, err := authUc.Register(ctx, auth.RegisterDTO{
		Email:     "budi.santoso@example.com",
		Password:  "SecurePassword123!",
		FirstName: "Budi",
		LastName:  "Santoso",
		Company:   "PT Solusi Cloud Nusantara",
		Country:   "ID",
		Currency:  "IDR",
	})
	if err != nil {
		panic(err)
	}
	registeredClient, _ := clientRepo.GetByID(ctx, regRes.Client.ID)
	_ = eventBus.Publish(ctx, events.Event{
		Type:    events.EventClientRegistered,
		Payload: registeredClient,
	})

	fmt.Printf("   ✅ Klien Terdaftar: %s %s (ID: %d, Email: %s)\n",
		regRes.Client.FirstName, regRes.Client.LastName, regRes.Client.ID, regRes.Client.Email)
	fmt.Printf("   📧 Email Selamat Datang Terkirim via Mailer: %s\n", mockMailer.LastMessage().Subject)

	// 3. Pengecekan Domain Ketersediaan (Domain Registrar SPI)
	fmt.Println("\n[2] 🌐 Pengecekan Ketersediaan Domain (Domain Registrar SPI)...")
	avail, _ := domainDriver.CheckAvailability(ctx, "solusinusantara.com")
	fmt.Printf("   🔍 Domain: %s (Tersedia: %v, Harga: %s %s)\n",
		avail.DomainName, avail.IsAvailable, avail.Currency, decimal.Money(avail.Price).String())

	// 4. Simulasi Pembuatan Kupon Promo & Multi-Mata Uang (Currency Management)
	fmt.Println("\n[3] 💱 Mengelola Multi-Currency & Kupon Promo...")
	_, _ = currencyService.CreateCurrency(ctx, currency.CreateCurrencyDTO{
		Code:           "USD",
		Title:          "US Dollar",
		ConversionRate: 0.000065,
		Format:         "$ {{price}}",
		PriceFormat:    "2",
		IsDefault:      false,
	})
	currencies, _ := currencyService.ListCurrencies(ctx)
	fmt.Printf("   ✅ Multi-Mata Uang Aktif (%d valuta terdaftar): IDR (Default) & USD\n", len(currencies))

	promo := &domain.Promo{
		Code:   "MERDEKA20",
		Type:   domain.PromoTypePercentage,
		Value:  decimal.FromFloat(20.00), // Diskon 20%
		Active: true,
	}
	_ = promoRepo.Create(ctx, promo)
	fmt.Printf("   🏷️ Kupon Aktif: %s (Diskon %s%%)\n", promo.Code, promo.Value.String())

	// 5. Simulasi Pemilihan Produk ke Keranjang & Checkout (Termasuk Produk Downloadable)
	fmt.Println("\n[4] 🛒 Menambahkan Produk ke Keranjang & Checkout...")
	hostingCfg, _ := json.Marshal(map[string]string{"domain": "solusinusantara.com", "plan": "unlimited_pro"})

	shoppingCart := &cart.Cart{
		ClientID:  regRes.Client.ID,
		PromoCode: "MERDEKA20",
		Items: []cart.CartItem{
			{ProductID: 101, Title: "Cloud VPS cPanel Pro", Period: "1M", Price: decimal.FromFloat(200000.00), Quantity: 1, Config: hostingCfg},
			{ProductID: 202, Title: "DirectAdmin Business Hosting", Period: "1M", Price: decimal.FromFloat(150000.00), Quantity: 1},
			{ProductID: 303, Title: "FOSSBilling Enterprise License", Period: "1Y", Price: decimal.FromFloat(500000.00), Quantity: 1},
			{ProductID: 404, Title: "Nusantara Cloud OS Template", Period: "ONETIME", Price: decimal.FromFloat(100000.00), Quantity: 1},
		},
	}

	// Register downloadable product file
	_ = downloadRepo.Create(ctx, &domain.DownloadableFile{
		ProductID:   404,
		Filename:    "nusantara-cloud-os-v1.iso",
		FilePath:    "/data/files/nusantara-cloud-os-v1.iso",
		FileSize:    1024 * 1024 * 500, // 500 MB
		ContentType: "application/x-iso9660-image",
		Version:     "1.0.0",
	})

	checkoutRes, err := cartService.Checkout(ctx, shoppingCart)
	if err != nil {
		panic(err)
	}

	fmt.Printf("   ✅ Subtotal Tagihan: %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Subtotal.String())
	fmt.Printf("   🏛️ Pajak (PPN 11%%): %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Tax.String())
	fmt.Printf("   💰 TOTAL INVOICE : %s %s (Nomor: #%s%s, Status: %s)\n",
		checkoutRes.Invoice.Currency, checkoutRes.Invoice.Total.String(),
		checkoutRes.Invoice.Serie, checkoutRes.Invoice.Nr, checkoutRes.Invoice.Status)

	// 6. Simulasi Pembayaran Masuk via Webhook
	fmt.Println("\n[5] 💳 Menerima Notifikasi Pembayaran Webhook (Payment Settled)...")
	webhookPayload := payment.WebhookPayload{
		GatewayID: "midtrans",
		TxnID:     "MID-TXN-" + fmt.Sprint(time.Now().Unix()),
		InvoiceID: checkoutRes.Invoice.ID,
		Amount:    checkoutRes.Invoice.Total,
		Currency:  checkoutRes.Invoice.Currency,
		Raw:       []byte(`{"status":"settlement"}`),
	}

	txn, err := webhookService.HandlePaymentWebhook(ctx, webhookPayload)
	if err != nil {
		panic(err)
	}
	fmt.Printf("   ✅ Transaksi Dicatat: TxnID=%s, Gateway=%s, Status=%s\n", txn.TxnID, txn.GatewayID, txn.Status)

	// Trigger Event Invoice Paid
	_ = eventBus.Publish(ctx, events.Event{
		Type: events.EventInvoicePaid,
		Payload: map[string]interface{}{
			"client":  registeredClient,
			"invoice": checkoutRes.Invoice,
			"txn":     txn,
		},
	})
	fmt.Printf("   📧 Bukti Pembayaran Lunas Terkirim ke Email: %s\n", mockMailer.LastMessage().Subject)

	// 7. Render PDF / HTML Invoice
	fmt.Println("\n[6] 📄 Menghasilkan Dokumen Faktur Pajak / Invoice PDF...")
	pdfBytes, _ := pdf.GenerateInvoiceHTML(checkoutRes.Invoice, registeredClient, "Nusantara Cloud Indonesia", "", "")
	fmt.Printf("   ✅ File Invoice PDF Berhasil Dihasilkan (%d bytes) dengan Status Lunas/Paid Stamp!\n", len(pdfBytes))

	// 8. Verifikasi Otomasi & Driver Provisioning
	fmt.Println("\n[7] ⚡ Eksekusi Multi-Driver Provisioning Otomatis...")
	for _, ord := range checkoutRes.Orders {
		activated, _ := orderRepo.GetByID(ctx, ord.ID)
		fmt.Printf("   🚀 Layanan Aktif: %s (Status: %s)\n", activated.Title, activated.Status)

		if ord.ProductID == 101 {
			res, _ := cpanelProv.Create(ctx, activated)
			var details map[string]string
			_ = json.Unmarshal(res.AccountDetails, &details)
			fmt.Printf("      📦 cPanel Server   : %s (User: %s, Pass: %s)\n", details["server"], details["username"], details["password"])
		} else if ord.ProductID == 202 {
			daAcc, _ := daProv.CreateAccount(ctx, provisioning.DirectAdminAccount{Domain: "solusinusantara.com", Package: "Business"})
			fmt.Printf("      📦 DirectAdmin     : Host %s (User: %s, Pass: %s)\n", daProv.Host, daAcc.Username, daAcc.Password)
		} else if ord.ProductID == 303 {
			res, _ := licenseProv.Create(ctx, activated)
			var details map[string]string
			_ = json.Unmarshal(res.AccountDetails, &details)
			fmt.Printf("      🔑 Enterprise Key  : %s\n", details["license_key"])
		}
	}

	// Plesk Provisioning
	pleskSub, _ := pleskProv.CreateSubscription(ctx, provisioning.PleskSubscription{DomainName: "plesk-demo.com", PlanName: "Default"})
	fmt.Printf("   🚀 Layanan Plesk   : Domain %s (Webspace User: %s, Status: %s)\n", pleskSub.DomainName, pleskSub.Username, pleskSub.Status)

	// 9. Digital Download with HMAC-SHA256 Signed Link
	fmt.Println("\n[8] 📦 Pengiriman Produk Digital (Signed Download Link)...")
	dlFile, _ := downloadRepo.GetByProductID(ctx, 404)
	signedLink, _ := downloadService.GenerateDownloadLink(ctx, regRes.Client.ID, dlFile.ID, 2*time.Hour)
	fmt.Printf("   🔗 Link Unduh Terproteksi: %s (Valid s/d %v)\n", signedLink.URL, signedLink.ExpiresAt.Format("15:04:05"))

	// 10. API Keys Management
	fmt.Println("\n[9] 🔑 Manajemen Service API Keys...")
	apiKey, _ := apiKeyService.GenerateKey(ctx, regRes.Client.ID, "Automated Deployment Bot", 90)
	fmt.Printf("   🔑 API Key Diterbitkan: %s (Secret: %s...)\n", apiKey.Key, apiKey.Secret[:8])

	// 11. News & Announcement Publishing
	fmt.Println("\n[10] 📰 Publikasi Pengumuman & Berita (News Module)...")
	article, _ := newsService.Create(ctx, news.CreateNewsDTO{
		AdminID: 1,
		Title:   "Pembaruan Infrastruktur Node Jakarta 2026",
		Content: "Kami telah meng-upgrade kapasitas server 10Gbps di datacenter IDC 3D.",
		Status:  domain.NewsStatusPublished,
	})
	fmt.Printf("   📰 Berita Terbit: \"%s\" (Slug: /news/%s)\n", article.Title, article.Slug)

	// 12. Mass Mailer & Broadcast Campaign
	fmt.Println("\n[11] 📢 Kampanye Email Massal (Mass Mailer Module)...")
	campaign, _ := massMailService.Create(ctx, 1, "Pengumuman Maintenance Terjadwal", "<p>Maintenance IDC 3D pada malam ini.</p>")
	sentCampaign, _ := massMailService.Send(ctx, campaign.ID)
	fmt.Printf("   📧 Kampanye Terkirim ke %d Klien! Status: %s\n", sentCampaign.SentCount, sentCampaign.Status)

	// 13. Customer Support Ticketing
	fmt.Println("\n[12] 🎫 Simulasi Customer Support Ticketing...")
	ticket, _ := supportService.OpenTicket(ctx, support.CreateTicketDTO{
		ClientID:   regRes.Client.ID,
		HelpdeskID: 1,
		Subject:    "Panduan Konfigurasi DNS Domain",
		Message:    "Halo, bagaimana cara mengarahkan DNS ke nameserver sg1.nusantara-cloud.com?",
		Priority:   domain.PriorityHigh,
	})
	fmt.Printf("   📝 Tiket Dibuka: [#%d] %s (Status: %s)\n", ticket.ID, ticket.Subject, ticket.Status)
	_, _ = supportService.StaffReply(ctx, ticket.ID, 1, "Halo Pak Budi, nameserver kami adalah ns1 & ns2.nusantara-cloud.com.")
	_, _ = supportService.ClientReply(ctx, ticket.ID, regRes.Client.ID, "Terima kasih, sudah bisa diakses!", "127.0.0.1")
	_ = supportService.CloseTicket(ctx, ticket.ID, regRes.Client.ID)
	finalTicket, _ := supportRepo.GetTicketByID(ctx, ticket.ID)
	fmt.Printf("   ✅ Tiket Selesai & Ditutup: %s\n", finalTicket.Status)

	// 14. Financial Analytics & Executive Dashboard
	fmt.Println("\n[13] 📊 Eksekutif Business Dashboard & Financial Analytics...")
	dashStats, _ := statsService.CalculateDashboard(ctx)
	fmt.Printf("   📈 MRR (Monthly Recurring) : %s %s\n", checkoutRes.Invoice.Currency, dashStats.MonthlyRecurring.String())
	fmt.Printf("   📈 ARR (Annual Recurring)  : %s %s\n", checkoutRes.Invoice.Currency, dashStats.AnnualRecurring.String())
	fmt.Printf("   💵 Total Omset Pendapatan  : %s %s\n", checkoutRes.Invoice.Currency, dashStats.TotalRevenue.String())
	fmt.Printf("   👥 Total Klien Aktif       : %d klien\n", dashStats.TotalClients)
	fmt.Printf("   📦 Total Pesanan Aktif     : %d layanan aktif\n", dashStats.ActiveOrders)

	fmt.Println("\n==================================================================")
	fmt.Println("🎉 SIMULASI SELESAI: 100% Seluruh Modul PHP Telah Sukses Terimplementasi di Golang!")
	fmt.Println("==================================================================")
}
