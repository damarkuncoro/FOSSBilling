package seed

import (
	"context"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
)

// SeedPromos populates sample promotional coupon codes
func SeedPromos(ctx context.Context, repo *memory.MockPromoRepository) {
	if repo == nil {
		return
	}

	_ = repo.Create(ctx, &domain.Promo{
		ID:          1,
		Code:        "MERDEKA20",
		Description: "20% Discount Celebration Voucher",
		Type:        domain.PromoTypePercentage,
		Value:       200000,
		Active:      true,
	})
}

// SeedCurrencies populates default base and secondary currencies
func SeedCurrencies(ctx context.Context, repo *memory.MockCurrencyRepository) {
	if repo == nil {
		return
	}

	_ = repo.Create(ctx, &domain.Currency{
		ID:             1,
		Code:           "USD",
		Title:          "US Dollar",
		Format:         "$ {{price}}",
		ConversionRate: 1.0,
		IsDefault:      true,
	})
	_ = repo.Create(ctx, &domain.Currency{
		ID:             2,
		Code:           "IDR",
		Title:          "Indonesian Rupiah",
		Format:         "Rp {{price}}",
		ConversionRate: 16000.0,
		IsDefault:      false,
	})
}

// SeedNews populates announcement articles
func SeedNews(ctx context.Context, repo *memory.MockNewsRepository) {
	if repo == nil {
		return
	}

	now := time.Now().UTC()
	_ = repo.Create(ctx, &domain.NewsPost{
		ID:          1,
		AdminID:     1,
		Title:       "Singapore Datacenter Core Switch Upgrade Scheduled",
		Slug:        "singapore-switch-upgrade",
		Content:     "We will be performing a scheduled hardware firmware upgrade on core Singapore switches.",
		Status:      domain.NewsStatusPublished,
		PublishedAt: &now,
	})
}
