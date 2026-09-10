package seed

import (
	"context"
	"time"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
)

func SeedPromos(ctx context.Context, r *memory.MockPromoRepository) {
	if r == nil { return }
	d := "20% Discount Celebration Voucher"
	_ = r.Create(ctx, &domain.Promo{ID: 1, Code: "MERDEKA20", Description: &d, Type: domain.PromoTypePercentage, Value: 200000, Active: true})
}

func SeedCurrencies(ctx context.Context, r *memory.MockCurrencyRepository) {
	if r == nil { return }
	_ = r.Create(ctx, &domain.Currency{ID: 1, Code: "USD", Title: "US Dollar", Format: "$ {{price}}", ConversionRate: 1.0, IsDefault: true})
	_ = r.Create(ctx, &domain.Currency{ID: 2, Code: "IDR", Title: "Indonesian Rupiah", Format: "Rp {{price}}", ConversionRate: 16000.0, IsDefault: false})
}

func SeedNews(ctx context.Context, r *memory.MockNewsRepository) {
	if r == nil { return }
	n := time.Now().UTC()
	_ = r.Create(ctx, &domain.NewsPost{ID: 1, AdminID: 1, Title: "Switch Upgrade", Slug: "singapore-switch-upgrade", Content: "Scheduled hardware upgrade.", Status: domain.NewsStatusPublished, PublishedAt: &n})
}
