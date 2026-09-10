package builder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type OrderBuilder struct {
	cid, pid, iid int64; tit, per, cur string; pri decimal.Money; st domain.OrderStatus; cfg map[string]any; days int
}

func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{per: "1M", cur: "USD", st: domain.OrderStatusPendingSetup, cfg: make(map[string]any), days: 30}
}

func (b *OrderBuilder) ForClient(id int64) *OrderBuilder { b.cid = id; return b }
func (b *OrderBuilder) ForProduct(id int64, t string) *OrderBuilder { b.pid, b.tit = id, t; return b }
func (b *OrderBuilder) WithInvoice(id int64) *OrderBuilder { b.iid = id; return b }
func (b *OrderBuilder) WithPeriod(p string) *OrderBuilder { b.per = p; return b }
func (b *OrderBuilder) WithPrice(p decimal.Money, c string) *OrderBuilder { b.pri, b.cur = p, c; return b }
func (b *OrderBuilder) WithConfig(k string, v any) *OrderBuilder { b.cfg[k] = v; return b }
func (b *OrderBuilder) AsActive() *OrderBuilder { b.st = domain.OrderStatusActive; return b }

func (b *OrderBuilder) Build() (*domain.Order, error) {
	if b.cid == 0 || b.pid == 0 || b.tit == "" { return nil, errors.New("missing fields") }
	cj, _ := json.Marshal(b.cfg); now := time.Now().UTC()
	var act, exp *time.Time
	if b.st == domain.OrderStatusActive { act = &now; e := now.AddDate(0, 0, b.days); exp = &e }
	var iid *int64; if b.iid > 0 { iid = &b.iid }
	return &domain.Order{ClientID: b.cid, ProductID: b.pid, InvoiceID: iid, Title: b.tit, Period: b.per, Price: b.pri, Currency: b.cur, Status: b.st, Config: cj, ActivatedAt: act, ExpiresAt: exp, NextDueDate: exp, CreatedAt: now, UpdatedAt: now}, nil
}
