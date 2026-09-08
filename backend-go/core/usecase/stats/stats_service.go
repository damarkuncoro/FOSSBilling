package stats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type RevenueTrend struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
	MRR     float64 `json:"mrr"`
}

type DashboardStats struct {
	TotalRevenue     decimal.Money  `json:"total_revenue"`
	MonthlyRecurring decimal.Money  `json:"mrr"`
	MonthlyRevenue   decimal.Money  `json:"monthly_revenue"` // Frontend compatibility
	AnnualRecurring  decimal.Money  `json:"arr"`
	TotalClients     int            `json:"total_clients"`
	ActiveClients    int            `json:"active_clients"` // Frontend compatibility
	TotalOrders      int            `json:"total_orders"`   // Frontend compatibility
	ActiveOrders     int            `json:"active_orders"`
	SuspendedOrders  int            `json:"suspended_orders"`
	PendingOrders    int            `json:"pending_orders"`
	UnpaidInvoices   int            `json:"unpaid_invoices"`
	PaidInvoices     int            `json:"paid_invoices"`
	OpenTickets      int            `json:"open_tickets"`
	ClosedTickets    int            `json:"closed_tickets"`
	RevenueTrends    []RevenueTrend `json:"revenue_trends"`
}

type FinancialReportSummary struct {
	MRR                 float64 `json:"mrr"`
	ARR                 float64 `json:"arr"`
	TotalRevenueMonth   float64 `json:"total_revenue_month"`
	TotalTaxCollected   float64 `json:"total_tax_collected"`
	ActiveSubscriptions int     `json:"active_subscriptions"`
	ChurnRate           float64 `json:"churn_rate"`
	MonthlyBreakdown    []struct {
		Month         string  `json:"month"`
		Revenue       float64 `json:"revenue"`
		Tax           float64 `json:"tax"`
		InvoicesCount int     `json:"invoices_count"`
	} `json:"monthly_breakdown"`
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
	stats := &DashboardStats{
		RevenueTrends: make([]RevenueTrend, 0),
	}

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
	monthlyRevMap := make(map[string]decimal.Money)
	if err == nil {
		var totalRevenue decimal.Money
		for _, inv := range invoices {
			if inv.Status == domain.InvoiceStatusPaid {
				stats.PaidInvoices++
				totalRevenue += inv.Total
				monthKey := inv.CreatedAt.Format("Jan")
				if inv.PaidAt != nil {
					monthKey = inv.PaidAt.Format("Jan")
				}
				monthlyRevMap[monthKey] += inv.Total
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

	// 5. Generate 6-month Revenue Trends
	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	currentMonth := time.Now().Month() // 1..12
	for i := 5; i >= 0; i-- {
		mIdx := (int(currentMonth) - 1 - i + 12) % 12
		mName := monthNames[mIdx]
		rev := monthlyRevMap[mName].ToFloat()
		if rev == 0 && i == 0 {
			rev = stats.TotalRevenue.ToFloat()
		}
		stats.RevenueTrends = append(stats.RevenueTrends, RevenueTrend{
			Month:   mName,
			Revenue: rev,
			MRR:     stats.MonthlyRecurring.ToFloat(),
		})
	}

	// 6. Compatibility mappings
	stats.MonthlyRevenue = stats.MonthlyRecurring
	stats.ActiveClients = stats.TotalClients
	stats.TotalOrders = stats.ActiveOrders + stats.SuspendedOrders + stats.PendingOrders

	return stats, nil
}

// GetFinancialReports provides in-depth fiscal analytics for the reporting module
func (s *StatsService) GetFinancialReports(ctx context.Context) (*FinancialReportSummary, error) {
	report := &FinancialReportSummary{
		MonthlyBreakdown: make([]struct {
			Month         string  `json:"month"`
			Revenue       float64 `json:"revenue"`
			Tax           float64 `json:"tax"`
			InvoicesCount int     `json:"invoices_count"`
		}, 0),
	}

	// 1. Calculate MRR/ARR and Active Subs from Orders
	orders, _, _ := s.orderRepo.List(ctx, 10000, 0)
	var mrr decimal.Money
	for _, o := range orders {
		if o.Status == domain.OrderStatusActive {
			report.ActiveSubscriptions++
			mrr += calculateMonthlyEquivalent(o.Price, o.Period)
		}
	}
	report.MRR = mrr.ToFloat()
	report.ARR = report.MRR * 12

	// 2. Aggregate Invoices by Month
	invoices, _, _ := s.invoiceRepo.List(ctx, 10000, 0)

	type monthStat struct {
		revenue  decimal.Money
		tax      decimal.Money
		invCount int
	}
	statsMap := make(map[string]*monthStat)
	monthKeys := make([]string, 0)

	now := time.Now().UTC()
	currentMonthKey := now.Format("Jan 2006")

	for _, inv := range invoices {
		if inv.Status != domain.InvoiceStatusPaid {
			continue
		}

		key := inv.CreatedAt.Format("Jan 2006")
		if inv.PaidAt != nil {
			key = inv.PaidAt.Format("Jan 2006")
		}

		if _, exists := statsMap[key]; !exists {
			statsMap[key] = &monthStat{}
			monthKeys = append(monthKeys, key)
		}

		statsMap[key].revenue += inv.Total
		statsMap[key].tax += inv.Tax
		statsMap[key].invCount++

		if key == currentMonthKey {
			report.TotalRevenueMonth = statsMap[key].revenue.ToFloat()
			report.TotalTaxCollected = statsMap[key].tax.ToFloat()
		}
	}

	// 3. Populate breakdown (last 5 months)
	for i := 4; i >= 0; i-- {
		d := now.AddDate(0, -i, 0)
		key := d.Format("Jan 2006")

		rev, tax, count := 0.0, 0.0, 0
		if ms, ok := statsMap[key]; ok {
			rev = ms.revenue.ToFloat()
			tax = ms.tax.ToFloat()
			count = ms.invCount
		}

		report.MonthlyBreakdown = append(report.MonthlyBreakdown, struct {
			Month         string  `json:"month"`
			Revenue       float64 `json:"revenue"`
			Tax           float64 `json:"tax"`
			InvoicesCount int     `json:"invoices_count"`
		}{
			Month:         key,
			Revenue:       rev,
			Tax:           tax,
			InvoicesCount: count,
		})
	}

	report.ChurnRate = 1.2 // Mock constant for now
	return report, nil
}

// GenerateInvoicesCSV returns a CSV string of all invoices for accounting
func (s *StatsService) GenerateInvoicesCSV(ctx context.Context) (string, error) {
	invoices, _, err := s.invoiceRepo.List(ctx, 50000, 0)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("Invoice ID,Number,Client ID,Status,Currency,Subtotal,Tax,Total,Tax Rate %,Due Date,Paid At,Created At\n")

	for _, inv := range invoices {
		paidAt := ""
		if inv.PaidAt != nil {
			paidAt = inv.PaidAt.Format("2006-01-02 15:04:05")
		}

		line := fmt.Sprintf("%d,%s%s,%d,%s,%s,%s,%s,%s,%.2f,%s,%s,%s\n",
			inv.ID, inv.Serie, inv.Nr, inv.ClientID, inv.Status, inv.Currency,
			inv.Subtotal.String(), inv.Tax.String(), inv.Total.String(),
			inv.TaxRate, inv.DueAt.Format("2006-01-02"),
			paidAt, inv.CreatedAt.Format("2006-01-02 15:04:05"),
		)
		sb.WriteString(line)
	}

	return sb.String(), nil
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
