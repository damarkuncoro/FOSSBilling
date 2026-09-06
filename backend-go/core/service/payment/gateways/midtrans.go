package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type MidtransGateway struct {
	serverKey string
	clientKey string
	isProd    bool
}

func NewMidtransGateway(serverKey, clientKey string, isProd bool) *MidtransGateway {
	return &MidtransGateway{
		serverKey: serverKey,
		clientKey: clientKey,
		isProd:    isProd,
	}
}

func (g *MidtransGateway) ID() string   { return "midtrans" }
func (g *MidtransGateway) Name() string { return "Midtrans Payment Gateway (QRIS, GoPay, VA)" }
func (g *MidtransGateway) Type() string { return "wallet" }

func (g *MidtransGateway) InitiatePayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResponse, error) {
	snapToken := fmt.Sprintf("SNAP-MID-%d-%d", req.InvoiceID, time.Now().Unix())
	redirectURL := fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/%s", snapToken)
	if g.isProd {
		redirectURL = fmt.Sprintf("https://app.midtrans.com/snap/v2/vtweb/%s", snapToken)
	}

	return &payment.PaymentResponse{
		GatewayID:     g.ID(),
		TransactionID: fmt.Sprintf("MID-ORDER-%d", req.InvoiceID),
		RedirectURL:   redirectURL,
		Token:         snapToken,
		Metadata: map[string]interface{}{
			"invoice_nr": req.InvoiceNr,
		},
	}, nil
}

func (g *MidtransGateway) ParseWebhook(r *http.Request) (*payment.WebhookResult, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var payload struct {
		TransactionStatus string `json:"transaction_status"`
		OrderID           string `json:"order_id"`
		GrossAmount       string `json:"gross_amount"`
		TransactionID     string `json:"transaction_id"`
		StatusCode        string `json:"status_code"`
		Currency          string `json:"currency"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// Extract invoice ID from order_id (e.g., "INV-101" or "101")
	invIDStr := strings.TrimPrefix(payload.OrderID, "INV-")
	invID, _ := strconv.ParseInt(invIDStr, 10, 64)

	amtFloat, _ := strconv.ParseFloat(payload.GrossAmount, 64)
	isPaid := payload.TransactionStatus == "settlement" || payload.TransactionStatus == "capture"

	return &payment.WebhookResult{
		GatewayID:     g.ID(),
		TransactionID: payload.TransactionID,
		InvoiceID:     invID,
		Amount:        decimal.FromFloat(amtFloat),
		Currency:      payload.Currency,
		IsPaid:        isPaid,
		RawPayload:    body,
	}, nil
}
