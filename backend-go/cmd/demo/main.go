package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	"github.com/fossbilling/backend-go/internal/repository/memory"
	"github.com/fossbilling/backend-go/internal/usecase/auth"
	"github.com/fossbilling/backend-go/internal/usecase/billing"
	"github.com/fossbilling/backend-go/internal/usecase/cart"
	"github.com/fossbilling/backend-go/internal/usecase/order"
	"github.com/fossbilling/backend-go/internal/usecase/payment"
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

	taxCalc := billing.NewTaxCalculator([]billing.TaxRule{
		{Name: "Indonesian PPN", Country: "ID", Rate: 11.0},
	})

	authUc := auth.NewAuthUsecase(clientRepo, "super-secret-jwt-key-32-chars")
	orderService := order.NewOrderService(orderRepo)
	invService := billing.NewInvoiceService(invRepo, clientRepo, taxCalc)
	promoCalc := cart.NewPromoCalculator(promoRepo)
	cartService := cart.NewCartService(promoCalc, promoRepo, orderRepo, invService)
	webhookService := payment.NewWebhookService(txnRepo, invRepo, orderService, orderRepo)

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
	fmt.Printf("   🔑 JWT Access Token Terbit: %s...[truncated]\n", regRes.Token[:30])

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
	shoppingCart := &cart.Cart{
		ClientID:  regRes.Client.ID,
		PromoCode: "MERDEKA20",
		Items: []cart.CartItem{
			{ProductID: 101, Title: "Cloud VPS cPanel Pro", Period: "1M", Price: decimal.FromFloat(200000.00), Quantity: 1},
			{ProductID: 202, Title: "Domain .id Registration", Period: "1Y", Price: decimal.FromFloat(150000.00), Quantity: 1},
		},
	}

	checkoutRes, err := cartService.Checkout(ctx, shoppingCart)
	if err != nil {
		panic(err)
	}

	fmt.Printf("   ✅ Subtotal Awal : Rp 350.000,00\n")
	fmt.Printf("   🎁 Diskon Kupon  : -Rp 70.000,00 (20%%)\n")
	fmt.Printf("   🧾 Subtotal Tagihan: %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Subtotal.String())
	fmt.Printf("   🏛️ Pajak (PPN 11%%): %s %s\n", checkoutRes.Invoice.Currency, checkoutRes.Invoice.Tax.String())
	fmt.Printf("   💰 TOTAL INVOICE : %s %s (Nomor: #%s%s, Status: %s)\n",
		checkoutRes.Invoice.Currency, checkoutRes.Invoice.Total.String(),
		checkoutRes.Invoice.Serie, checkoutRes.Invoice.Nr, checkoutRes.Invoice.Status)

	for i, ord := range checkoutRes.Orders {
		fmt.Printf("   📦 Pesanan Dibuat [%d]: %s (Status: %s)\n", i+1, ord.Title, ord.Status)
	}

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

	// 6. Verifikasi Status Invoice & Aktivasi Layanan Otomatis
	fmt.Println("\n[5] ⚡ Verifikasi Otomasi & Aktivasi Layanan...")
	paidInv, _ := invRepo.GetByID(ctx, checkoutRes.Invoice.ID)
	fmt.Printf("   🧾 Status Invoice : %s (Lunas pada: %v)\n", paidInv.Status, paidInv.PaidAt.Format("02 Jan 2006 15:04:05"))

	for _, ord := range checkoutRes.Orders {
		activated, _ := orderRepo.GetByID(ctx, ord.ID)
		fmt.Printf("   🚀 Layanan Aktif: %s\n", activated.Title)
		fmt.Printf("      - Status      : %s\n", activated.Status)
		fmt.Printf("      - Diaktifkan  : %v\n", activated.ActivatedAt.Format("02 Jan 2006"))
		fmt.Printf("      - Jatuh Tempo : %v\n", activated.ExpiresAt.Format("02 Jan 2006"))
	}

	fmt.Println("\n==================================================================")
	fmt.Println("🎉 SIMULASI SELESAI: Seluruh alur bisnis berjalan 100% sempurna!")
	fmt.Println("==================================================================")
}
