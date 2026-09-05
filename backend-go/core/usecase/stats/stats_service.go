package stats

import (
	"context"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type DashboardStats struct {
	TotalRevenue      decimal.Money `json:"total_revenue"`
	MonthlyRecurring  decimal.Money `json:"mrr"`
	AnnualRecurring   decimal.Money `json:"arr"`
	TotalClients      int           `json:"total_clients"`
	ActiveOrders      int           `json:"active_orders"`
	SuspendedOrders   int           `json:"suspended_orders"`
	PendingOrders     int           `json:"pending_orders"`
	UnpaidInvoices    int           `json:"unpaid_invoices"`
	PaidInvoices      int           `json:"paid_invoices"`
	OpenTickets       int           `json:"open_tickets"`
	ClosedTickets     int           `json:"closed_tickets"`
}

type StatsService struct {
	clientRepo  domain.ClientRepository
	orderRepo   domain.OrderRepository
	invoiceRepo domain.InvoiceRepository
	supportRepo domain.SupportRepository
}

func NewStatsService(
	clientRepo domain.ClientRepository,
	orderRepo domain.OrderRepository,
	invoiceRepo domain.InvoiceRepository,
	supportRepo domain.SupportRepository,
) *StatsService {
	return &StatsService{
		clientRepo:  clientRepo,
		orderRepo:   orderRepo,
		invoiceRepo: invoiceRepo,
		supportRepo: supportRepo,
	}
}

// CalculateDashboard aggregates real-time business and operations KPIs
func (s *StatsService) CalculateDashboard(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 1. Client count
	_, totalClients, err := s.clientRepo.List(ctx, 1, 0)
	if err == nil {
		stats.TotalClients = totalClients
	}

	// 2. Orders & MRR calculation
	orders, _, err := s.orderRepo.List(ctx, 10000, 0)
	if err == nil {
		var mrrTotal decimal.Money
		for _, o := range orders {
			switch o.Status {
			case domain.OrderStatusActive:
				stats.ActiveOrders++
				// Normalize recurring period into monthly equivalent
				mrrTotal += calculateMonthlyEquivalent(o.Price, o.Period)
			case domain.OrderStatusSuspended:
				stats.SuspendedOrders++
			case domain.OrderStatusPendingSetup:
				stats.PendingOrders++
			}
		}
		stats.MonthlyRecurring = mrrTotal
		stats.AnnualRecurring = mrrTotal * 12
	}

	// 3. Invoices & Total Collected Revenue
	invoices, _, err := s.invoiceRepo.List(ctx, 10000, 0)
	if err == nil {
		var totalRevenue decimal.Money
		for _, inv := range invoices {
			if inv.Status == domain.InvoiceStatusPaid {
				stats.PaidInvoices++
				totalRevenue += inv.Total
			} else if inv.Status == domain.InvoiceStatusUnpaid {
				stats.UnpaidInvoices++
			}
		}
		stats.TotalRevenue = totalRevenue
	}

	// 4. Support Tickets
	tickets, _, err := s.supportRepo.ListTickets(ctx, 10000, 0)
	if err == nil {
		for _, t := range tickets {
			if t.Status == domain.TicketStatusClosed {
				stats.ClosedTickets++
			} else {
				stats.OpenTickets++
			}
		}
	}

	return stats, nil
}

// calculateMonthlyEquivalent converts different billing periods to monthly amounts
func calculateMonthlyEquivalent(price decimal.Money, periodStr string) decimal.Money {
	periodStr = strings.ToUpper(strings.TrimSpace(periodStr))
	switch periodStr {
	case "1W":
		return price * 4 // Approx 4 weeks in a month
	case "2W":
		return price * 2
	case "1M":
		return price
	case "3M":
		return price / 3
	case "6M":
		return price / 6
	case "1Y":
		return price / 12
	case "2Y":
		return price / 24
	case "3Y":
		return price / 36
	default:
		return 0 // Onetime / Free do not contribute to recurring MRR
	}
}
