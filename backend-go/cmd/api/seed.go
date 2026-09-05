package main

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func seedMemoryData(
	ctx context.Context,
	clientRepo *memory.MockClientRepository,
	staffRepo *memory.MockStaffRepository,
	promoRepo *memory.MockPromoRepository,
	newsRepo *memory.MockNewsRepository,
	currencyRepo *memory.MockCurrencyRepository,
	invoiceRepo *memory.MockInvoiceRepository,
) {
	passHash, _ := auth.HashPassword("Password123!")

	// Seed Client
	_ = clientRepo.Create(ctx, &domain.Client{
		ID:           1,
		Email:        "client@fossbilling.org",
		PasswordHash: passHash,
		FirstName:    "Budi",
		LastName:     "Santoso",
		Company:      "PT Solusi Cloud Nusantara",
		Country:      "ID",
		Currency:     "USD",
		Status:       domain.ClientStatusActive,
	})

	_ = clientRepo.AddBalanceTransaction(ctx, &domain.ClientBalance{
		ClientID:    1,
		Type:        "credit",
		Amount:      5000000,
		Description: "Initial wallet deposit balance",
	})

	// Seed Invoices for Client 1
	now := time.Now().UTC()
	dueDate := now.Add(14 * 24 * time.Hour)
	paidDate := now.Add(-2 * 24 * time.Hour)

	_ = invoiceRepo.Create(ctx, &domain.Invoice{
		ID:           1,
		Serie:        "INV",
		Nr:           "2026-0001",
		ClientID:     1,
		Status:       domain.InvoiceStatusPaid,
		Currency:     "USD",
		CurrencyRate: 1.0,
		Subtotal:     decimal.FromFloat(100.0),
		Tax:          decimal.FromFloat(11.0),
		Total:        decimal.FromFloat(111.0),
		TaxRate:      11.0,
		DueAt:        dueDate,
		PaidAt:       &paidDate,
		CreatedAt:    now.Add(-3 * 24 * time.Hour),
		UpdatedAt:    paidDate,
	}, []domain.InvoiceItem{
		{
			ID:        1,
			InvoiceID: 1,
			Title:     "Cloud VPS Premium (2 vCPU, 4GB RAM, 80GB NVMe)",
			Price:     decimal.FromFloat(100.0),
			Quantity:  1,
			Taxable:   true,
			CreatedAt: now.Add(-3 * 24 * time.Hour),
		},
	})

	_ = invoiceRepo.Create(ctx, &domain.Invoice{
		ID:           2,
		Serie:        "INV",
		Nr:           "2026-0002",
		ClientID:     1,
		Status:       domain.InvoiceStatusUnpaid,
		Currency:     "USD",
		CurrencyRate: 1.0,
		Subtotal:     decimal.FromFloat(49.0),
		Tax:          decimal.FromFloat(5.39),
		Total:        decimal.FromFloat(54.39),
		TaxRate:      11.0,
		DueAt:        dueDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, []domain.InvoiceItem{
		{
			ID:        2,
			InvoiceID: 2,
			Title:     "cPanel Web Hosting Professional - 1 Year Renewal",
			Price:     decimal.FromFloat(49.0),
			Quantity:  1,
			Taxable:   true,
			CreatedAt: now,
		},
	})

	// Seed Admin Group & Staff
	group := &domain.AdminGroup{
		ID:   1,
		Name: "Super Administrators",
		Permissions: map[string][]string{
			"clients": {"*"}, "orders": {"*"}, "support": {"*"}, "system": {"*"},
			"billing": {"*"}, "staff": {"*"}, "news": {"*"}, "currencies": {"*"},
		},
	}
	_ = staffRepo.CreateGroup(ctx, group)

	adminPass, _ := auth.HashPassword("SuperSecretAdmin123!")
	_ = staffRepo.Create(ctx, &domain.Staff{
		ID:           1,
		GroupID:      group.ID,
		Email:        "admin@fossbilling.org",
		PasswordHash: adminPass,
		Name:         "Super Administrator",
		Role:         domain.StaffRoleSuperAdmin,
		Status:       "active",
	})

	// Seed Promo
	_ = promoRepo.Create(ctx, &domain.Promo{
		ID:          1,
		Code:        "MERDEKA20",
		Description: "20% Discount Celebration Voucher",
		Type:        domain.PromoTypePercentage,
		Value:       200000,
		Active:      true,
	})

	// Seed Currencies
	_ = currencyRepo.Create(ctx, &domain.Currency{
		ID: 1, Code: "USD", Title: "US Dollar", Format: "$ {{price}}", ConversionRate: 1.0, IsDefault: true,
	})
	_ = currencyRepo.Create(ctx, &domain.Currency{
		ID: 2, Code: "IDR", Title: "Indonesian Rupiah", Format: "Rp {{price}}", ConversionRate: 16000.0, IsDefault: false,
	})

	// Seed News
	_ = newsRepo.Create(ctx, &domain.NewsPost{
		ID:          1,
		AdminID:     1,
		Title:       "Singapore Datacenter Core Switch Upgrade Scheduled",
		Slug:        "singapore-switch-upgrade",
		Content:     "We will be performing a scheduled hardware firmware upgrade on core Singapore switches.",
		Status:      domain.NewsStatusPublished,
		PublishedAt: &now,
	})
}
