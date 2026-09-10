package payment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
)

var (
	ErrDuplicateTransaction = errors.New("transaction already processed")
	ErrGatewayNotFound      = errors.New("gateway not found")
)

type WebhookPayload struct {
	GatewayID, TxnID, Currency string; InvoiceID int64; Amount decimal.Money; Raw []byte
}

type WebhookService struct {
	tr domain.TransactionRepository; ir domain.InvoiceRepository; reg *payment.GatewayRegistry; eb *events.EventBus
}

func NewWebhookService(tr domain.TransactionRepository, ir domain.InvoiceRepository, r *payment.GatewayRegistry, eb ...*events.EventBus) *WebhookService {
	var bus *events.EventBus; if len(eb) > 0 { bus = eb[0] }
	return &WebhookService{tr, ir, r, bus}
}

func (s *WebhookService) ProcessRawWebhook(r *http.Request, gid string) (*payment.WebhookResult, error) {
	gw, err := s.reg.Get(gid); if err != nil { return nil, ErrGatewayNotFound }; return gw.ParseWebhook(r)
}

func (s *WebhookService) HandlePaymentWebhook(ctx context.Context, p WebhookPayload) (*domain.Transaction, error) {
	if ex, _ := s.tr.GetByTxnID(ctx, p.GatewayID, p.TxnID); ex != nil && ex.Status == domain.TransactionStatusComplete { return ex, ErrDuplicateTransaction }
	inv, err := s.ir.GetByID(ctx, p.InvoiceID); if err != nil { return nil, err }
	if p.Amount == 0 { p.Amount = inv.Total }; if p.Currency == "" { p.Currency = inv.Currency }

	raw := json.RawMessage("{}"); if len(p.Raw) > 0 && json.Valid(p.Raw) { raw = json.RawMessage(p.Raw) }
	txn := &domain.Transaction{InvoiceID: &inv.ID, GatewayID: p.GatewayID, TxnID: p.TxnID, Type: domain.TransactionTypePayment, Amount: p.Amount, Currency: p.Currency, Status: domain.TransactionStatusComplete, RawPayload: raw}
	if err := s.tr.Create(ctx, txn); err != nil { return nil, err }

	now := time.Now().UTC()
	if inv.Status != domain.InvoiceStatusPaid { _ = s.ir.MarkAsPaid(ctx, inv.ID, now) }
	if s.eb != nil { _ = s.eb.Publish(ctx, events.Event{Type: events.EventInvoicePaid, Payload: domain.InvoicePaidPayload{InvoiceID: inv.ID, ClientID: inv.ClientID, Amount: p.Amount, Currency: p.Currency, GatewayID: p.GatewayID, TxnID: p.TxnID, PaidAt: now}}) }
	return txn, nil
}
