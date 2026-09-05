package cart

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

var (
	ErrPromoExpired       = errors.New("promo code has expired")
	ErrPromoNotStarted    = errors.New("promo code is not yet active")
	ErrPromoInactive      = errors.New("promo code is disabled")
	ErrPromoMaxUses       = errors.New("promo code maximum usage limit reached")
	ErrPromoAlreadyUsed   = errors.New("promo code can only be used once per client")
)

type PromoCalculator struct {
	promoRepo domain.PromoRepository
}

func NewPromoCalculator(promoRepo domain.PromoRepository) *PromoCalculator {
	return &PromoCalculator{promoRepo: promoRepo}
}

// ValidatePromo validates coupon eligibility for a client
func (c *PromoCalculator) ValidatePromo(ctx context.Context, promo *domain.Promo, clientID int64, now time.Time) error {
	if !promo.Active {
		return ErrPromoInactive
	}

	if promo.StartDate != nil && now.Before(*promo.StartDate) {
		return ErrPromoNotStarted
	}

	if promo.EndDate != nil && now.After(*promo.EndDate) {
		return ErrPromoExpired
	}

	if promo.MaxUses > 0 && promo.UsedCount >= promo.MaxUses {
		return ErrPromoMaxUses
	}

	if promo.OncePerClient && clientID > 0 {
		count, err := c.promoRepo.GetRedemptionCount(ctx, promo.ID, clientID)
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrPromoAlreadyUsed
		}
	}

	return nil
}

// CalculateDiscount calculates the discount value against a subtotal
func (c *PromoCalculator) CalculateDiscount(subtotal decimal.Money, promo *domain.Promo) decimal.Money {
	if subtotal <= 0 || promo == nil {
		return 0
	}

	if promo.Type == domain.PromoTypePercentage {
		// e.g. Value = 20.00% (stored as Money: 200000)
		percentage := promo.Value.ToFloat()
		if percentage <= 0 {
			return 0
		}
		discountFloat := subtotal.ToFloat() * (percentage / 100.0)
		discount := decimal.FromFloat(discountFloat)
		if discount > subtotal {
			return subtotal
		}
		return discount
	}

	// Absolute fixed discount
	if promo.Value > subtotal {
		return subtotal
	}
	return promo.Value
}
