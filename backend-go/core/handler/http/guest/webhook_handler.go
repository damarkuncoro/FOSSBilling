package guest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type WebhookHandler struct {
	webhookService *payment.WebhookService
}

func NewWebhookHandler(webhookService *payment.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService}
}

func (h *WebhookHandler) HandleGatewayWebhook(w http.ResponseWriter, r *http.Request) {
	// 1. Determine which gateway this webhook is for
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	gatewayID := ""
	for i, part := range pathParts {
		if (part == "gateways" || part == "webhook") && i+1 < len(pathParts) {
			gatewayID = pathParts[i+1]
			break
		}
	}

	if gatewayID == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_GATEWAY", "Gateway ID not specified in path", nil)
		return
	}

	// 2. Delegate parsing and SIGNATURE VERIFICATION to the specific gateway driver
	// This is CRITICAL for security (BUG-32 FIX)
	result, err := h.webhookService.ProcessRawWebhook(r, gatewayID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "WEBHOOK_VERIFICATION_FAILED", err.Error(), nil)
		return
	}

	if !result.IsPaid {
		response.JSON(w, http.StatusOK, map[string]string{"status": "ignored", "reason": "event_not_paid"}, nil)
		return
	}

	// 3. Process the verified payment in the system
	payload := payment.WebhookPayload{
		GatewayID: gatewayID,
		TxnID:     result.TransactionID,
		InvoiceID: result.InvoiceID,
		Amount:    result.Amount,
		Currency:  result.Currency,
		Raw:       result.RawPayload,
	}

	txn, err := h.webhookService.HandlePaymentWebhook(r.Context(), payload)
	if err != nil {
		if errors.Is(err, payment.ErrDuplicateTransaction) {
			response.JSON(w, http.StatusOK, map[string]interface{}{
				"status":  "duplicate_ignored",
				"message": "Transaction already processed",
				"txn_id":  payload.TxnID,
			}, nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "PROCESSING_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":         "processed",
		"transaction_id": txn.ID,
		"txn_id":         txn.TxnID,
		"invoice_id":     txn.InvoiceID,
	}, nil)
}
