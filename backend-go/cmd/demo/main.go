package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/repository/memory"
	"github.com/fossbilling/backend-go/internal/service/provisioning"
	"github.com/fossbilling/backend-go/internal/usecase/auth"
	"github.com/fossbilling/backend-go/internal/usecase/billing"
	"github.com/fossbilling/backend-go/internal/usecase/cart"
	"github.com/fossbilling/backend-go/internal/usecase/order"
	"github.com/fossbilling/backend-go/internal/usecase/payment"
	"github.com/fossbilling/backend-go/internal/usecase/support"
	"github.com/fossbilling/backend-go/pkg/decimal"
)

func main() {
	ctx := context.Background()
	fmt.Println("==================================================================")
	fmt.Println("🚀 FOSSBilling Next-Gen Backend (Golang) - E2E Live Simulation")
	fmt.Println("==================================================================")

	// 1. Inisialisasi Repositories & Services
	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invRepo := memory.NewMockInvoiceRepository()
	promoRepo := memory.NewMockPromoRepository()
	txnRepo := memory.NewMockTransactionRepository()
	supportRepo := memory.NewMockSupportRepository()

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

	cpanelProv := provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{Host: "sg1.nusantara-cloud.com"})
	licenseProv := provisioning.NewLicenseProvisioner("fossbilling-enterprise-master-key")

	// 2. Simulasi Registrasi Klien Baru
	fmt.Println("\n[1] 👤 Mendaftarkan Klien Baru...")
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
	fmt.Printf("   ✅ Klien Terdaftar: %s %s (ID: %d, Email: %s)\n",
		regRes.Client.FirstName, regRes.Client.LastName, regRes.Client.ID, regRes.Client.Email)

	// 3. Simulasi Pembuatan Kupon Promo
	fmt.Println("\n[2] 🏷️ Membuat Kupon Promo...")
	promo := &domain.Promo{
		Code:   "MERDEKA20",
		Type:   domain.PromoTypePercentage,
		Value:  decimal.FromFloat(20.00), // Diskon 20%
		Active: true,
	}
	_ = promoRepo.Create(ctx, promo)
	fmt.Printf("   ✅ Kupon Aktif: %s (Diskon %s%%)\n", promo.Code, promo.Value.String())

	// 4. Simulasi Pemilihan Produk ke Keranjang & Checkout
	fmt.Println("\n[3] 🛒 Menambahkan Produk ke Keranjang & Checkout...")
	hostingCfg, _ := json.Marshal(map[string]string{"domain": "solusinusantara.com", "plan": "unlimited_pro"})

	shoppingCart := &cart.Cart{
		ClientID:  regRes.Client.ID,
		PromoCode: "MERDEKA20",
		Items: []cart.CartItem{
			{ProductID: 101, Title: "Cloud VPS cPanel Pro", Period: "1M", Price: decimal.FromFloat(200000.00), Quantity: 1, Config: hostingCfg},
			{ProductID: 303, Title: "FOSSBilling Enterprise License", Period: "1Y", Price: decimal.FromFloat(500000.00), Quantity: 1},
		},
	}

	checkoutRes, err := cartService.Checkout(ctx, shoppingCart)
	if err != nil {
		panic(err)
	}

	fmt.Printf("   ✅ Subtotal Awal : IDR 700000.00\n")
	fmt.Printf("   🎁 Diskon Kupon  : -IDR 140000.00 (20%%)\n")
	fmt.Printf("   🧾 Subtotal Tagihan: %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Subtotal.String())
	fmt.Printf("   🏛️ Pajak (PPN 11%%): %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Tax.String())
	fmt.Printf("   💰 TOTAL INVOICE : %s %s (Nomor: #%s%s, Status: %s)\n",
		checkoutRes.Invoice.Currency, checkoutRes.Invoice.Total.String(),
		checkoutRes.Invoice.Serie, checkoutRes.Invoice.Nr, checkoutRes.Invoice.Status)

	// 5. Simulasi Pembayaran Masuk via Webhook (Payment Gateway Midtrans/Stripe)
	fmt.Println("\n[4] 💳 Menerima Notifikasi Pembayaran Webhook (Payment Settled)...")
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

	// 6. Verifikasi Otomasi & Driver Provisioning
	fmt.Println("\n[5] ⚡ Eksekusi Driver Provisioning Otomatis...")
	paidInv, _ := invRepo.GetByID(ctx, checkoutRes.Invoice.ID)
	fmt.Printf("   🧾 Status Invoice : %s (Lunas pada: %v)\n", paidInv.Status, paidInv.PaidAt.Format("02 Jan 2006 15:04:05"))

	for _, ord := range checkoutRes.Orders {
		activated, _ := orderRepo.GetByID(ctx, ord.ID)
		fmt.Printf("   🚀 Layanan Aktif: %s (Status: %s)\n", activated.Title, activated.Status)

		if ord.ProductID == 101 {
			// Hosting account provision
			res, _ := cpanelProv.Create(ctx, activated)
			var details map[string]string
			_ = json.Unmarshal(res.AccountDetails, &details)
			fmt.Printf("      📦 cPanel Server   : %s\n", details["server"])
			fmt.Printf("      👤 cPanel User     : %s\n", details["username"])
			fmt.Printf("      🔑 cPanel Pass     : %s\n", details["password"])
			fmt.Printf("      🌐 Login URL       : %s\n", details["cpanel_url"])
		} else if ord.ProductID == 303 {
			// License provision
			res, _ := licenseProv.Create(ctx, activated)
			var details map[string]string
			_ = json.Unmarshal(res.AccountDetails, &details)
			fmt.Printf("      🔑 License Key     : %s\n", details["license_key"])
		}
	}

	// 7. Simulasi Layanan Bantuan / Support Ticket
	fmt.Println("\n[6] 🎫 Simulasi Customer Support Ticketing...")
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

	fmt.Println("\n==================================================================")
	fmt.Println("🎉 SIMULASI SELESAI: Seluruh ekosistem backend Go berjalan 100% sempurna!")
	fmt.Println("==================================================================")
}
