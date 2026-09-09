package payment

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"net/http"
)

var (
	ErrDuplicateTransaction = errors.New("transaction has already been processed")
	ErrGatewayNotFound      = errors.New("payment gateway driver not found")
)

type WebhookPayload struct {
	GatewayID string        `json:"gateway_id"`
	TxnID     string        `json:"txn_id"`
	InvoiceID int64         `json:"invoice_id"`
	Amount    decimal.Money `json:"amount"`
	Currency  string        `json:"currency"`
	Raw       []byte        `json:"-"`
}

type WebhookService struct {
	txnRepo         domain.TransactionRepository
	invoiceRepo     domain.InvoiceRepository
	gatewayRegistry *payment.GatewayRegistry
	eventBus        *events.EventBus
}

func NewWebhookService(
	txnRepo domain.TransactionRepository,
	invoiceRepo domain.InvoiceRepository,
	gatewayRegistry *payment.GatewayRegistry,
	eventBus ...*events.EventBus,
) *WebhookService {
	var bus *events.EventBus
	if len(eventBus) > 0 {
		bus = eventBus[0]
	}
	return &WebhookService{
		txnRepo:         txnRepo,
		invoiceRepo:     invoiceRepo,
		gatewayRegistry: gatewayRegistry,
		eventBus:        bus,
	}
}

// ProcessRawWebhook identifies the gateway and calls its driver to verify signature and parse data
func (s *WebhookService) ProcessRawWebhook(r *http.Request, gatewayID string) (*payment.WebhookResult, error) {
	gw, err := s.gatewayRegistry.Get(gatewayID)
	if err != nil {
		return nil, ErrGatewayNotFound
	}

	return gw.ParseWebhook(r)
}

// HandlePaymentWebhook processes inbound payment IPN/webhooks idempotently and triggers order activation via events
func (s *WebhookService) HandlePaymentWebhook(ctx context.Context, payload WebhookPayload) (*domain.Transaction, error) {
	// 1. Idempotency Check
	existing, err := s.txnRepo.GetByTxnID(ctx, payload.GatewayID, payload.TxnID)
	if err == nil && existing != nil {
		if existing.Status == domain.TransactionStatusComplete {
			return existing, ErrDuplicateTransaction
		}
	}

	// 2. Fetch Invoice
	invoice, err := s.invoiceRepo.GetByID(ctx, payload.InvoiceID)
	if err != nil {
		return nil, err
	}

	if payload.Amount == 0 {
		payload.Amount = invoice.Total
	}
	if payload.Currency == "" {
		payload.Currency = invoice.Currency
	}

	// 3. Record Transaction
	var raw json.RawMessage
	if len(payload.Raw) > 0 && json.Valid(payload.Raw) {
		raw = json.RawMessage(payload.Raw)
	} else {
		raw = json.RawMessage("{}")
	}

	txn := &domain.Transaction{
		InvoiceID:  &invoice.ID,
		GatewayID:  payload.GatewayID,
		TxnID:      payload.TxnID,
		Type:       domain.TransactionTypePayment,
		Amount:     payload.Amount,
		Currency:   payload.Currency,
		Status:     domain.TransactionStatusComplete,
		RawPayload: raw,
	}

	if err := s.txnRepo.Create(ctx, txn); err != nil {
		return nil, err
	}

	// 4. Mark Invoice as Paid
	now := time.Now().UTC()
	if invoice.Status != domain.InvoiceStatusPaid {
		if err := s.invoiceRepo.MarkAsPaid(ctx, invoice.ID, now); err != nil {
			return nil, err
		}
	}

	// 5. Publish Event (Listeners will handle order activation and provisioning)
	if s.eventBus != nil {
		err := s.eventBus.Publish(ctx, events.Event{
			Type: events.EventInvoicePaid,
			Payload: domain.InvoicePaidPayload{
				InvoiceID: invoice.ID,
				ClientID:  invoice.ClientID,
				Amount:    payload.Amount,
				Currency:  payload.Currency,
				GatewayID: payload.GatewayID,
				TxnID:     payload.TxnID,
				PaidAt:    now,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return txn, nil
}
