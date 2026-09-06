package payment_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment/gateways"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func TestPaymentGateways_Registry(t *testing.T) {
	registry := payment.NewGatewayRegistry()

	midtrans := gateways.NewMidtransGateway("server_key", "client_key", false)
	stripe := gateways.NewStripeGateway("sk_test_123", "pk_test_123", "whsec_123")
	paypal := gateways.NewPayPalGateway("client_id", "secret", false)
	bankTransfer := gateways.NewBankTransferGateway("BCA", "1234567890", "PT FOSSBilling")

	registry.Register(midtrans)
	registry.Register(stripe)
	registry.Register(paypal)
	registry.Register(bankTransfer)

	list := registry.List()
	if len(list) != 4 {
		t.Fatalf("expected 4 gateways, got %d", len(list))
	}

	gw, err := registry.Get("stripe")
	if err != nil || gw.ID() != "stripe" {
		t.Fatalf("failed to retrieve stripe gateway: %v", err)
	}

	_, err = registry.Get("nonexistent")
	if err == nil {
		t.Fatalf("expected error for nonexistent gateway")
	}
}

func TestPaymentGateways_InitiateAndWebhook(t *testing.T) {
	ctx := context.Background()
	req := payment.PaymentRequest{
		InvoiceID:   101,
		InvoiceNr:   "INV-2026-0101",
		Amount:      decimal.FromFloat(150.00),
		Currency:    "USD",
		ClientEmail: "client@fossbilling.org",
		ClientName:  "Test Client",
	}

	// 1. Stripe
	stripe := gateways.NewStripeGateway("sk_test", "pk_test", "whsec_test")
	resp, err := stripe.InitiatePayment(ctx, req)
	if err != nil || resp.Token == "" {
		t.Fatalf("Stripe InitiatePayment failed: %v", err)
	}

	stripeWebhookBody := `{"type":"checkout.session.completed","data":{"object":{"id":"cs_test_101_999","amount_total":15000,"currency":"usd","payment_status":"paid","metadata":{"invoice_id":"101"}}}}`
	httpReq, _ := http.NewRequest("POST", "/webhook/stripe", bytes.NewBufferString(stripeWebhookBody))
	whRes, err := stripe.ParseWebhook(httpReq)
	if err != nil || !whRes.IsPaid || whRes.InvoiceID != 101 {
		t.Fatalf("Stripe ParseWebhook failed: %v, result: %+v", err, whRes)
	}

	// 2. PayPal
	paypal := gateways.NewPayPalGateway("client_id", "secret", false)
	pResp, err := paypal.InitiatePayment(ctx, req)
	if err != nil || pResp.Token == "" {
		t.Fatalf("PayPal InitiatePayment failed: %v", err)
	}

	paypalWebhookBody := `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"id":"PAYPAL-101-CAPTURE","custom_id":"101","status":"COMPLETED","amount":{"value":"150.00","currency_code":"USD"}}}`
	httpReq, _ = http.NewRequest("POST", "/webhook/paypal", bytes.NewBufferString(paypalWebhookBody))
	pWhRes, err := paypal.ParseWebhook(httpReq)
	if err != nil || !pWhRes.IsPaid || pWhRes.InvoiceID != 101 {
		t.Fatalf("PayPal ParseWebhook failed: %v, result: %+v", err, pWhRes)
	}

	// 3. Midtrans
	midtrans := gateways.NewMidtransGateway("server_key", "client_key", false)
	mResp, err := midtrans.InitiatePayment(ctx, req)
	if err != nil || mResp.Token == "" {
		t.Fatalf("Midtrans InitiatePayment failed: %v", err)
	}

	midtransWebhookBody := `{"order_id":"INV-101","transaction_status":"settlement","gross_amount":"150.00","transaction_id":"txn-mid-101","status_code":"200","currency":"IDR"}`
	httpReq, _ = http.NewRequest("POST", "/webhook/midtrans", bytes.NewBufferString(midtransWebhookBody))
	mWhRes, err := midtrans.ParseWebhook(httpReq)
	if err != nil || !mWhRes.IsPaid || mWhRes.InvoiceID != 101 {
		t.Fatalf("Midtrans ParseWebhook failed: %v, result: %+v", err, mWhRes)
	}

	// 4. Bank Transfer
	bt := gateways.NewBankTransferGateway("Mandiri", "1400011223344", "PT Solusi Cloud")
	btResp, err := bt.InitiatePayment(ctx, req)
	if err != nil || btResp.TransactionID == "" {
		t.Fatalf("BankTransfer InitiatePayment failed: %v", err)
	}
}
