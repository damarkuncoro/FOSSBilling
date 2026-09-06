package spec

import (
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type UnpaidInvoicesSpec struct{}

func NewUnpaidInvoicesSpec() Specification[*domain.Invoice] {
	return &UnpaidInvoicesSpec{}
}

func (s *UnpaidInvoicesSpec) IsSatisfiedBy(inv *domain.Invoice) bool {
	if inv == nil {
		return false
	}
	return inv.Status == domain.InvoiceStatusUnpaid
}

type OverdueInvoicesSpec struct {
	asOf time.Time
}

func NewOverdueInvoicesSpec(asOf time.Time) Specification[*domain.Invoice] {
	return &OverdueInvoicesSpec{asOf: asOf}
}

func (s *OverdueInvoicesSpec) IsSatisfiedBy(inv *domain.Invoice) bool {
	if inv == nil {
		return false
	}
	return inv.Status == domain.InvoiceStatusUnpaid && inv.DueAt.Before(s.asOf)
}
