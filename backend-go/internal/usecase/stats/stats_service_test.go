package stats_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/stats"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func TestStatsService_CalculateDashboard(t *testing.T) {
	ctx := context.Background()

	clientRepo := memory.NewMockClientRepository()
	orderRepo := memory.NewMockOrderRepository()
	invoiceRepo := memory.NewMockInvoiceRepository()
	supportRepo := memory.NewMockSupportRepository()

	statsService := stats.NewStatsService(clientRepo, orderRepo, invoiceRepo, supportRepo)

	// 1. Add Clients
	_ = clientRepo.Create(ctx, &domain.Client{Email: "c1@example.com"})
	_ = clientRepo.Create(ctx, &domain.Client{Email: "c2@example.com"})

	// 2. Add Orders
	// Order 1: 1M at $10.00 (100000)
	_ = orderRepo.Create(ctx, &domain.Order{
		ClientID: 1,
		Status:   domain.OrderStatusActive,
		Period:   "1M",
		Price:    100000,
	})
	// Order 2: 1Y at $120.00 (1200000) -> MRR: $10.00 (100000)
	_ = orderRepo.Create(ctx, &domain.Order{
		ClientID: 2,
		Status:   domain.OrderStatusActive,
		Period:   "1Y",
		Price:    1200000,
	})
	// Order 3: Suspended
	_ = orderRepo.Create(ctx, &domain.Order{
		ClientID: 1,
		Status:   domain.OrderStatusSuspended,
		Period:   "1M",
		Price:    50000,
	})

	// 3. Add Invoices
	_ = invoiceRepo.Create(ctx, &domain.Invoice{
		ClientID: 1,
		Status:   domain.InvoiceStatusPaid,
		Total:    100000,
	}, nil)
	_ = invoiceRepo.Create(ctx, &domain.Invoice{
		ClientID: 2,
		Status:   domain.InvoiceStatusPaid,
		Total:    1200000,
	}, nil)
	_ = invoiceRepo.Create(ctx, &domain.Invoice{
		ClientID: 1,
		Status:   domain.InvoiceStatusUnpaid,
		Total:    50000,
	}, nil)

	// 4. Add Tickets
	_ = supportRepo.CreateTicket(ctx, &domain.Ticket{
		ClientID: 1,
		Subject:  "Open issue",
		Status:   domain.TicketStatusOpen,
	}, nil)
	_ = supportRepo.CreateTicket(ctx, &domain.Ticket{
		ClientID: 2,
		Subject:  "Closed issue",
		Status:   domain.TicketStatusClosed,
	}, nil)

	// Calculate Dashboard
	d, err := statsService.CalculateDashboard(ctx)
	if err != nil {
		t.Fatalf("CalculateDashboard failed: %v", err)
	}

	if d.TotalClients != 2 {
		t.Errorf("TotalClients = %d; want 2", d.TotalClients)
	}
	if d.ActiveOrders != 2 {
		t.Errorf("ActiveOrders = %d; want 2", d.ActiveOrders)
	}
	if d.SuspendedOrders != 1 {
		t.Errorf("SuspendedOrders = %d; want 1", d.SuspendedOrders)
	}
	// MRR should be 100000 (1M) + 100000 (1Y / 12) = 200000 ($20.00)
	expectedMRR := decimal.Money(200000)
	if d.MonthlyRecurring != expectedMRR {
		t.Errorf("MRR = %d; want %d", d.MonthlyRecurring, expectedMRR)
	}
	// ARR should be 200000 * 12 = 2400000 ($240.00)
	expectedARR := decimal.Money(2400000)
	if d.AnnualRecurring != expectedARR {
		t.Errorf("ARR = %d; want %d", d.AnnualRecurring, expectedARR)
	}
	// Total Revenue Paid = 100000 + 1200000 = 1300000 ($130.00)
	expectedRev := decimal.Money(1300000)
	if d.TotalRevenue != expectedRev {
		t.Errorf("TotalRevenue = %d; want %d", d.TotalRevenue, expectedRev)
	}
	if d.PaidInvoices != 2 || d.UnpaidInvoices != 1 {
		t.Errorf("PaidInvoices = %d, UnpaidInvoices = %d; want 2, 1", d.PaidInvoices, d.UnpaidInvoices)
	}
	if d.OpenTickets != 1 || d.ClosedTickets != 1 {
		t.Errorf("OpenTickets = %d, ClosedTickets = %d; want 1, 1", d.OpenTickets, d.ClosedTickets)
	}
}
