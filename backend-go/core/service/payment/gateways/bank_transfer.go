package gateways

import (
	"context"
	"fmt"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
)

type BankTransferGateway struct {
	bankName      string
	accountNumber string
	accountHolder string
}

func NewBankTransferGateway(bankName, accountNumber, accountHolder string) *BankTransferGateway {
	return &BankTransferGateway{
		bankName:      bankName,
		accountNumber: accountNumber,
		accountHolder: accountHolder,
	}
}

func (g *BankTransferGateway) ID() string   { return "bank_transfer" }
func (g *BankTransferGateway) Name() string { return "Bank Transfer (Manual Confirmation)" }
func (g *BankTransferGateway) Type() string { return "manual" }

func (g *BankTransferGateway) InitiatePayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResponse, error) {
	txnID := fmt.Sprintf("BT-MANUAL-%d", req.InvoiceID)

	return &payment.PaymentResponse{
		GatewayID:     g.ID(),
		TransactionID: txnID,
		Token:         txnID,
		Metadata: map[string]interface{}{
			"invoice_nr":     req.InvoiceNr,
			"bank_name":      g.bankName,
			"account_number": g.accountNumber,
			"account_holder": g.accountHolder,
			"instructions":   fmt.Sprintf("Silakan transfer ke rekening %s: %s a/n %s dengan menyertakan referensi invoice %s", g.bankName, g.accountNumber, g.accountHolder, req.InvoiceNr),
		},
	}, nil
}

func (g *BankTransferGateway) ParseWebhook(r *http.Request) (*payment.WebhookResult, error) {
	// Manual bank transfer does not have automatic webhook callbacks; handled via admin approval
	return nil, fmt.Errorf("manual bank transfer does not support automated webhook callbacks")
}
