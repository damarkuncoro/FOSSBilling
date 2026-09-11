package affiliate

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type AffiliateService struct {
	affRepo    domain.AffiliateRepository
	clientRepo domain.ClientRepository
	orderRepo  domain.OrderRepository
}

func NewAffiliateService(ar domain.AffiliateRepository, cr domain.ClientRepository, or domain.OrderRepository) *AffiliateService {
	return &AffiliateService{affRepo: ar, clientRepo: cr, orderRepo: or}
}

func (s *AffiliateService) GetAffiliate(ctx context.Context, clientID int64) (*domain.Affiliate, error) {
	return s.affRepo.GetByClientID(ctx, clientID)
}

func (s *AffiliateService) ActivateAffiliate(ctx context.Context, clientID int64) (*domain.Affiliate, error) {
	aff, err := s.affRepo.GetByClientID(ctx, clientID)
	if err == nil {
		return aff, nil
	}

	newAff := &domain.Affiliate{
		ClientID:       clientID,
		CommissionRate: 10.0, // Default 10%
		Status:         domain.AffiliateStatusActive,
		Balance:        0,
		TotalEarned:    0,
	}
	if err := s.affRepo.Create(ctx, newAff); err != nil {
		return nil, err
	}
	return newAff, nil
}

func (s *AffiliateService) ProcessCommission(ctx context.Context, orderID int64) error {
	o, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	c, err := s.clientRepo.GetByID(ctx, o.ClientID)
	if err != nil || c.ReferrerID == nil {
		return nil // No referrer, nothing to do
	}

	aff, err := s.affRepo.GetByClientID(ctx, *c.ReferrerID)
	if err != nil || aff.Status != domain.AffiliateStatusActive {
		return nil
	}

	// Calculate commission: Price * (Rate / 100)
	commissionAmount := decimal.Money(o.Price.ToFloat() * (aff.CommissionRate / 100.0))

	ref := &domain.AffiliateReferral{
		AffiliateID: aff.ID,
		ClientID:    c.ID,
		OrderID:     o.ID,
		Amount:      commissionAmount,
		Status:      "approved",
	}

	if err := s.affRepo.AddReferral(ctx, ref); err != nil {
		return err
	}

	// Update affiliate balance
	aff.Balance += commissionAmount
	aff.TotalEarned += commissionAmount
	return s.affRepo.Update(ctx, aff)
}
