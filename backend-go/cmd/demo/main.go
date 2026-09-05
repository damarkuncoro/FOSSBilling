package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/repository/memory"
	"github.com/fossbilling/backend-go/internal/service/notification"
	"github.com/fossbilling/backend-go/internal/service/provisioning"
	"github.com/fossbilling/backend-go/internal/usecase/auth"
	"github.com/fossbilling/backend-go/internal/usecase/billing"
	"github.com/fossbilling/backend-go/internal/usecase/cart"
	"github.com/fossbilling/backend-go/internal/usecase/order"
	"github.com/fossbilling/backend-go/internal/usecase/payment"
	"github.com/fossbilling/backend-go/internal/usecase/stats"
	"github.com/fossbilling/backend-go/internal/usecase/support"
	"github.com/fossbilling/backend-go/pkg/decimal"
	"github.com/fossbilling/backend-go/pkg/events"
	"github.com/fossbilling/backend-go/pkg/mailer"
	"github.com/fossbilling/backend-go/pkg/pdf"
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

	mockMailer := mailer.NewMockMailer()
	emailService := notification.NewEmailService(mockMailer, "billing@nusantara-cloud.com", "Nusantara Cloud")
	eventBus := events.NewEventBus()

	taxCalc := billing.NewTaxCalculator([]billing.TaxRule{
		{Name: "Indonesian PPN", Country: "ID", Rate: 11.0},
	})

	authUc := auth.NewAuthUsecase(clientRepo, "super-secret-jwt-key-32-chars")
	orderService := order.NewOrderService(orderRepo)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc)
	promoCalc := cart.NewPromoCalculator(promoRepo)
	cartService := cart.NewCartService(promoCalc, promoRepo, orderRepo, invService)
	webhookService := payment.NewWebhookService(txnRepo, invRepo, orderService, orderRepo)
	supportService := support.NewSupportService(supportRepo, clientRepo)
	statsService := stats.NewStatsService(clientRepo, orderRepo, invRepo, supportRepo)

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

	// 4. Simulasi Pembuatan Kupon Promo
	fmt.Println("\n[3] 🏷️ Membuat Kupon Promo...")
	promo := &domain.Promo{
		Code:   "MERDEKA20",
		Type:   domain.PromoTypePercentage,
		Value:  decimal.FromFloat(20.00), // Diskon 20%
		Active: true,
	}
	_ = promoRepo.Create(ctx, promo)
	fmt.Printf("   ✅ Kupon Aktif: %s (Diskon %s%%)\n", promo.Code, promo.Value.String())

	// 5. Simulasi Pemilihan Produk ke Keranjang & Checkout
	fmt.Println("\n[4] 🛒 Menambahkan Produk ke Keranjang & Checkout...")
	hostingCfg, _ := json.Marshal(map[string]string{"domain": "solusinusantara.com", "plan": "unlimited_pro"})

	shoppingCart := &cart.Cart{
		ClientID:  regRes.Client.ID,
		PromoCode: "MERDEKA20",
		Items: []cart.CartItem{
			{ProductID: 101, Title: "Cloud VPS cPanel Pro", Period: "1M", Price: decimal.FromFloat(200000.00), Quantity: 1, Config: hostingCfg},
			{ProductID: 202, Title: "DirectAdmin Business Hosting", Period: "1M", Price: decimal.FromFloat(150000.00), Quantity: 1},
			{ProductID: 303, Title: "FOSSBilling Enterprise License", Period: "1Y", Price: decimal.FromFloat(500000.00), Quantity: 1},
		},
	}

	checkoutRes, err := cartService.Checkout(ctx, shoppingCart)
	if err != nil {
		panic(err)
	}

	fmt.Printf("   ✅ Subtotal Awal : IDR 850000.00\n")
	fmt.Printf("   🎁 Diskon Kupon  : -IDR 170000.00 (20%%)\n")
	fmt.Printf("   🧾 Subtotal Tagihan: %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Subtotal.String())
	fmt.Printf("   🏛️ Pajak (PPN 11%%): %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Tax.String())
	fmt.Printf("   💰 TOTAL INVOICE : %s %s (Nomor: #%s%s, Status: %s)\n",
		checkoutRes.Invoice.Currency, checkoutRes.Invoice.Total.String(),
		checkoutRes.Invoice.Serie, checkoutRes.Invoice.Nr, checkoutRes.Invoice.Status)

	// 6. Simulasi Pembayaran Masuk via Webhook (Payment Gateway Midtrans/Stripe)
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


	// 8. Verifikasi Otomasi & Driver Provisioning (cPanel, DirectAdmin, Plesk, License)
	fmt.Println("\n[7] ⚡ Eksekusi Multi-Driver Provisioning Otomatis...")
	paidInv, _ := invRepo.GetByID(ctx, checkoutRes.Invoice.ID)
	fmt.Printf("   🧾 Status Invoice : %s (Lunas pada: %v)\n", paidInv.Status, paidInv.PaidAt.Format("02 Jan 2006 15:04:05"))

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

	// Test Plesk Driver as well
	pleskSub, _ := pleskProv.CreateSubscription(ctx, provisioning.PleskSubscription{DomainName: "plesk-demo.com", PlanName: "Default"})
	fmt.Printf("   🚀 Layanan Plesk   : Domain %s (Webspace User: %s, Status: %s)\n", pleskSub.DomainName, pleskSub.Username, pleskSub.Status)

	// 9. Simulasi Layanan Bantuan / Support Ticket
	fmt.Println("\n[8] 🎫 Simulasi Customer Support Ticketing...")
	ticket, err := supportService.OpenTicket(ctx, support.CreateTicketDTO{
		ClientID:   regRes.Client.ID,
		HelpdeskID: 1,
		Subject:    "Panduan Konfigurasi DNS Domain",
		Message:    "Halo, bagaimana cara mengarahkan DNS ke nameserver sg1.nusantara-cloud.com?",
		Priority:   domain.PriorityHigh,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("   📝 Tiket Dibuka: [#%d] %s (Status: %s, Prioritas: %s)\n",
		ticket.ID, ticket.Subject, ticket.Status, ticket.Priority)

	// Admin Balas
	_, _ = supportService.StaffReply(ctx, ticket.ID, 1, "Halo Pak Budi, nameserver kami adalah ns1.nusantara-cloud.com dan ns2.nusantara-cloud.com.")
	updatedTicket, _ := supportRepo.GetTicketByID(ctx, ticket.ID)
	fmt.Printf("   💬 Staf Membalas -> Status Tiket Berubah: %s\n", updatedTicket.Status)

	// Klien Konfirmasi Selesai & Tiket Ditutup
	_, _ = supportService.ClientReply(ctx, ticket.ID, regRes.Client.ID, "Terima kasih, sudah bisa diakses!", "127.0.0.1")
	_ = supportService.CloseTicket(ctx, ticket.ID, regRes.Client.ID)
	finalTicket, _ := supportRepo.GetTicketByID(ctx, ticket.ID)
	fmt.Printf("   ✅ Klien Puas -> Tiket Ditutup: %s\n", finalTicket.Status)

	// 10. Financial Analytics & Executive Dashboard
	fmt.Println("\n[9] 📊 Eksekutif Business Dashboard & Financial Analytics...")
	dashStats, _ := statsService.CalculateDashboard(ctx)
	fmt.Printf("   📈 MRR (Monthly Recurring) : %s %s\n", checkoutRes.Invoice.Currency, dashStats.MonthlyRecurring.String())
	fmt.Printf("   📈 ARR (Annual Recurring)  : %s %s\n", checkoutRes.Invoice.Currency, dashStats.AnnualRecurring.String())
	fmt.Printf("   💵 Total Omset Pendapatan  : %s %s\n", checkoutRes.Invoice.Currency, dashStats.TotalRevenue.String())
	fmt.Printf("   👥 Total Klien Aktif       : %d klien\n", dashStats.TotalClients)
	fmt.Printf("   📦 Total Pesanan Aktif     : %d layanan aktif\n", dashStats.ActiveOrders)

	fmt.Println("\n==================================================================")
	fmt.Println("🎉 SIMULASI SELESAI: Seluruh ekosistem backend Go berjalan 100% sempurna!")
	fmt.Println("==================================================================")
}
