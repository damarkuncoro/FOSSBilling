package payment

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/order"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

var (
	ErrDuplicateTransaction = errors.New("transaction has already been processed")
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
	txnRepo      domain.TransactionRepository
	invoiceRepo  domain.InvoiceRepository
	orderService *order.OrderService
	orderRepo    domain.OrderRepository
}

func NewWebhookService(
	txnRepo domain.TransactionRepository,
	invoiceRepo domain.InvoiceRepository,
	orderService *order.OrderService,
	orderRepo domain.OrderRepository,
) *WebhookService {
	return &WebhookService{
		txnRepo:      txnRepo,
		invoiceRepo:  invoiceRepo,
		orderService: orderService,
		orderRepo:    orderRepo,
	}
}

// HandlePaymentWebhook processes inbound payment IPN/webhooks idempotently and triggers order activation
func (s *WebhookService) HandlePaymentWebhook(ctx context.Context, payload WebhookPayload) (*domain.Transaction, error) {
	// 1. Idempotency Check: Prevent duplicate transaction processing
	existing, err := s.txnRepo.GetByTxnID(ctx, payload.GatewayID, payload.TxnID)
	if err == nil && existing != nil {
		if existing.Status == domain.TransactionStatusComplete {
			return existing, ErrDuplicateTransaction
		}
	}

	// 2. Fetch target Invoice
	invoice, err := s.invoiceRepo.GetByID(ctx, payload.InvoiceID)
	if err != nil {
		return nil, err
	}

	// 3. Record Transaction
	txn := &domain.Transaction{
		InvoiceID:  &invoice.ID,
		GatewayID:  payload.GatewayID,
		TxnID:      payload.TxnID,
		Type:       domain.TransactionTypePayment,
		Amount:     payload.Amount,
		Currency:   payload.Currency,
		Status:     domain.TransactionStatusComplete,
		RawPayload: json.RawMessage(payload.Raw),
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

	// 5. Automated Service Provisioning / Order Activation
	for _, it := range invoice.Items {
		if it.OrderID != nil {
			ord, err := s.orderRepo.GetByID(ctx, *it.OrderID)
			if err == nil && ord != nil {
				if ord.Status == domain.OrderStatusPendingSetup || ord.Status == domain.OrderStatusSuspended {
					_, _ = s.orderService.Activate(ctx, ord.ID, now)
				} else if ord.Status == domain.OrderStatusActive {
					_, _ = s.orderService.Renew(ctx, ord.ID)
				}
			}
		}
	}

	return txn, nil
}
