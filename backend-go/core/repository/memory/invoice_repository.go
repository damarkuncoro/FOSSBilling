package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockInvoiceRepository struct {
	mu           sync.RWMutex
	invoices     map[int64]*domain.Invoice
	invoiceItems map[int64][]domain.InvoiceItem
	nextID       int64
	nextItemID   int64
}

func NewMockInvoiceRepository() *MockInvoiceRepository {
	return &MockInvoiceRepository{
		invoices:     make(map[int64]*domain.Invoice),
		invoiceItems: make(map[int64][]domain.InvoiceItem),
		nextID:       1,
		nextItemID:   1,
	}
}

func (r *MockInvoiceRepository) GetByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	inv, ok := r.invoices[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *inv
	cp.Items = r.invoiceItems[id]
	return &cp, nil
}

func (r *MockInvoiceRepository) ListByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Invoice, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []*domain.Invoice
	for _, inv := range r.invoices {
		if inv.ClientID == clientID {
			cp := *inv
			cp.Items = r.invoiceItems[inv.ID]
			matched = append(matched, &cp)
		}
	}

	total := len(matched)
	if offset >= total {
		return []*domain.Invoice{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return matched[offset:end], total, nil
}

func (r *MockInvoiceRepository) List(ctx context.Context, limit, offset int) ([]*domain.Invoice, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []*domain.Invoice
	for _, inv := range r.invoices {
		cp := *inv
		cp.Items = r.invoiceItems[inv.ID]
		all = append(all, &cp)
	}

	total := len(all)
	if offset >= total {
		return []*domain.Invoice{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}


func (r *MockInvoiceRepository) Create(ctx context.Context, inv *domain.Invoice, items []domain.InvoiceItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	inv.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	inv.CreatedAt = now
	inv.UpdatedAt = now
	if inv.Status == "" {
		inv.Status = domain.InvoiceStatusUnpaid
	}
	if inv.Nr == "" {
		inv.Nr = fmt.Sprintf("%05d", inv.ID)
	}

	var savedItems []domain.InvoiceItem
	for _, it := range items {
		it.ID = r.nextItemID
		r.nextItemID++
		it.InvoiceID = inv.ID
		it.CreatedAt = now
		savedItems = append(savedItems, it)
	}

	cp := *inv
	r.invoices[inv.ID] = &cp
	r.invoiceItems[inv.ID] = savedItems
	return nil
}

func (r *MockInvoiceRepository) MarkAsPaid(ctx context.Context, id int64, paidAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	inv, ok := r.invoices[id]
	if !ok {
		return appErrors.ErrNotFound
	}

	inv.Status = domain.InvoiceStatusPaid
	inv.PaidAt = &paidAt
	inv.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *MockInvoiceRepository) Update(ctx context.Context, inv *domain.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.invoices[inv.ID]; !ok {
		return appErrors.ErrNotFound
	}
	inv.UpdatedAt = time.Now().UTC()
	cp := *inv
	r.invoices[inv.ID] = &cp
	return nil
}
