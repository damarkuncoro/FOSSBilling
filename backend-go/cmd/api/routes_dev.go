package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

// registerDevRoutes adds helper endpoints for manual testing in non-production environments
func registerDevRoutes(mux *http.ServeMux, h *AppHandlers) {
	// Endpoint to simulate a payment webhook for any invoice
	mux.HandleFunc("POST /api/v1/dev/simulate-payment/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		invID, _ := strconv.ParseInt(idStr, 10, 64)

		payload := payment.WebhookPayload{
			GatewayID: "midtrans",
			TxnID:     "DEV-SIM-" + idStr + "-" + strconv.FormatInt(time.Now().Unix(), 10),
			InvoiceID: invID,
			Amount:    0,
			Currency:  "USD",
		}

		err := h.GuestWebhook.HandleSimulation(r.Context(), payload)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "SIM_FAILED", err.Error(), nil)
			return
		}

		response.JSON(w, http.StatusOK, map[string]string{
			"message": "Payment simulation successful for Invoice #" + idStr + ". Orders will be provisioned by background worker within 1 minute.",
		}, nil)
	})
}
