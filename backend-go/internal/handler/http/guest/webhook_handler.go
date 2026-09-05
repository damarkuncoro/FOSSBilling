package guest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type WebhookHandler struct {
	webhookService *payment.WebhookService
}

func NewWebhookHandler(webhookService *payment.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService}
}

func (h *WebhookHandler) HandleGatewayWebhook(w http.ResponseWriter, r *http.Request) {
	// Extract gateway name from path, e.g., /api/v1/guest/gateways/stripe/webhook
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	gatewayID := "generic"
	for i, part := range pathParts {
		if part == "gateways" && i+1 < len(pathParts) {
			gatewayID = pathParts[i+1]
			break
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Failed to read request payload", nil)
		return
	}

	var payload payment.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid webhook JSON payload", nil)
		return
	}

	if payload.GatewayID == "" {
		payload.GatewayID = gatewayID
	}
	payload.Raw = body

	txn, err := h.webhookService.HandlePaymentWebhook(r.Context(), payload)
	if err != nil {
		if errors.Is(err, payment.ErrDuplicateTransaction) {
			response.JSON(w, http.StatusOK, map[string]interface{}{
				"status":  "duplicate_ignored",
				"message": "Transaction already processed",
				"txn_id":  txn.TxnID,
			}, nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "WEBHOOK_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":         "processed",
		"transaction_id": txn.ID,
		"txn_id":         txn.TxnID,
		"invoice_id":     txn.InvoiceID,
	}, nil)
}
