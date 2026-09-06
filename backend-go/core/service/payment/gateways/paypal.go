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

type PayPalGateway struct {
	clientID string
	secret   string
	isProd   bool
}

func NewPayPalGateway(clientID, secret string, isProd bool) *PayPalGateway {
	return &PayPalGateway{
		clientID: clientID,
		secret:   secret,
		isProd:   isProd,
	}
}

func (g *PayPalGateway) ID() string   { return "paypal" }
func (g *PayPalGateway) Name() string { return "PayPal Express Checkout" }
func (g *PayPalGateway) Type() string { return "wallet" }

func (g *PayPalGateway) InitiatePayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResponse, error) {
	orderID := fmt.Sprintf("PAYPAL-ORD-%d-%d", req.InvoiceID, time.Now().Unix())
	baseURL := "https://www.sandbox.paypal.com/checkoutnow"
	if g.isProd {
		baseURL = "https://www.paypal.com/checkoutnow"
	}
	redirectURL := fmt.Sprintf("%s?token=%s", baseURL, orderID)

	return &payment.PaymentResponse{
		GatewayID:     g.ID(),
		TransactionID: orderID,
		RedirectURL:   redirectURL,
		Token:         orderID,
		Metadata: map[string]interface{}{
			"invoice_nr": req.InvoiceNr,
			"client_id":  g.clientID,
		},
	}, nil
}

func (g *PayPalGateway) ParseWebhook(r *http.Request) (*payment.WebhookResult, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var payload struct {
		ID           string `json:"id"`
		EventType    string `json:"event_type"`
		Resource     struct {
			ID          string `json:"id"`
			Status      string `json:"status"`
			CustomID    string `json:"custom_id"`
			Amount      struct {
				Total    string `json:"total"`
				Value    string `json:"value"`
				Currency string `json:"currency_code"`
			} `json:"amount"`
		} `json:"resource"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	invID, _ := strconv.ParseInt(payload.Resource.CustomID, 10, 64)
	if invID == 0 {
		parts := strings.Split(payload.Resource.ID, "-")
		if len(parts) >= 3 {
			invID, _ = strconv.ParseInt(parts[2], 10, 64)
		}
	}

	valStr := payload.Resource.Amount.Value
	if valStr == "" {
		valStr = payload.Resource.Amount.Total
	}
	amtFloat, _ := strconv.ParseFloat(valStr, 64)
	isPaid := payload.EventType == "PAYMENT.CAPTURE.COMPLETED" || payload.Resource.Status == "COMPLETED"

	return &payment.WebhookResult{
		GatewayID:     g.ID(),
		TransactionID: payload.Resource.ID,
		InvoiceID:     invID,
		Amount:        decimal.FromFloat(amtFloat),
		Currency:      strings.ToUpper(payload.Resource.Amount.Currency),
		IsPaid:        isPaid,
		RawPayload:    body,
	}, nil
}
