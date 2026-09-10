package builder

import (
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type InvoiceBuilder struct {
	cid int64; ser, nr, cur string; crate, trate float64; due int; pr *domain.Promo; its []domain.InvoiceItem
}

func NewInvoiceBuilder() *InvoiceBuilder {
	return &InvoiceBuilder{ser: "INV", cur: "USD", crate: 1.0, due: 14, its: []domain.InvoiceItem{}}
}

func (b *InvoiceBuilder) ForClient(id int64) *InvoiceBuilder { b.cid = id; return b }
func (b *InvoiceBuilder) WithSerieAndNr(s, n string) *InvoiceBuilder { b.ser, b.nr = s, n; return b }
func (b *InvoiceBuilder) WithCurrency(c string, r float64) *InvoiceBuilder { b.cur, b.crate = c, r; return b }
func (b *InvoiceBuilder) WithDueDays(d int) *InvoiceBuilder { b.due = d; return b }
func (b *InvoiceBuilder) WithTaxRate(r float64) *InvoiceBuilder { b.trate = r; return b }
func (b *InvoiceBuilder) ApplyPromo(p *domain.Promo) *InvoiceBuilder { b.pr = p; return b }
func (b *InvoiceBuilder) AddItem(t string, p decimal.Money, q int, tx bool) *InvoiceBuilder {
	b.its = append(b.its, domain.InvoiceItem{Title: t, Price: p, Quantity: q, Taxable: tx}); return b
}

func (b *InvoiceBuilder) Build() (*domain.Invoice, []domain.InvoiceItem, error) {
	if b.cid == 0 || len(b.its) == 0 { return nil, nil, errors.New("missing fields") }
	var sub, txSub decimal.Money
	for _, it := range b.its { lt := it.Price * decimal.Money(it.Quantity); sub += lt; if it.Taxable { txSub += lt } }
	var dis decimal.Money
	if b.pr != nil && b.pr.Active {
		if b.pr.Type == domain.PromoTypePercentage { dis = decimal.FromFloat(sub.ToFloat() * (float64(b.pr.Value) / 1000000.0)) } else { dis = decimal.FromFloat(float64(b.pr.Value) / 10000.0) }
		if dis > sub { dis = sub }
	}
	tax := decimal.FromFloat(txSub.ToFloat() * (b.trate / 100.0))
	now := time.Now().UTC(); if b.nr == "" { b.nr = fmt.Sprintf("%d-%04d", now.Year(), now.Unix()%10000) }
	return &domain.Invoice{Serie: b.ser, Nr: b.nr, ClientID: b.cid, Status: domain.InvoiceStatusUnpaid, Currency: b.cur, CurrencyRate: b.crate, Subtotal: sub, Tax: tax, Total: sub - dis + tax, TaxRate: b.trate, DueAt: now.AddDate(0, 0, b.due), CreatedAt: now, UpdatedAt: now}, b.its, nil
}
