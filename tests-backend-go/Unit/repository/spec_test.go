package repository_test

import (
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/spec"
)

func TestSpecificationPattern_Invoices(t *testing.T) {
	now := time.Now().UTC()
	pastDue := now.Add(-48 * time.Hour)
	futureDue := now.Add(48 * time.Hour)

	invoices := []*domain.Invoice{
		{ID: 1, Status: domain.InvoiceStatusPaid, DueAt: pastDue},
		{ID: 2, Status: domain.InvoiceStatusUnpaid, DueAt: pastDue},
		{ID: 3, Status: domain.InvoiceStatusUnpaid, DueAt: futureDue},
	}

	// 1. Unpaid spec
	unpaidSpec := spec.NewUnpaidInvoicesSpec()
	unpaid := spec.Filter(invoices, unpaidSpec)
	if len(unpaid) != 2 {
		t.Fatalf("expected 2 unpaid invoices, got %d", len(unpaid))
	}

	// 2. Overdue spec
	overdueSpec := spec.NewOverdueInvoicesSpec(now)
	overdue := spec.Filter(invoices, overdueSpec)
	if len(overdue) != 1 || overdue[0].ID != 2 {
		t.Fatalf("expected invoice #2 to be overdue, got %+v", overdue)
	}
}

func TestSpecificationPattern_Orders(t *testing.T) {
	orders := []*domain.Order{
		{ID: 101, ClientID: 1, Status: domain.OrderStatusActive},
		{ID: 102, ClientID: 1, Status: domain.OrderStatusSuspended},
		{ID: 103, ClientID: 2, Status: domain.OrderStatusActive},
	}

	// 1. Client 1 Active Orders (And combination)
	client1Spec := spec.NewClientOrdersSpec(1)
	activeSpec := spec.NewActiveOrdersSpec()
	client1ActiveSpec := spec.And(client1Spec, activeSpec)

	matches := spec.Filter(orders, client1ActiveSpec)
	if len(matches) != 1 || matches[0].ID != 101 {
		t.Fatalf("expected order #101, got %+v", matches)
	}
}
