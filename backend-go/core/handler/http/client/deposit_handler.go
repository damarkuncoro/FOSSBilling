package client

import (
	"encoding/json"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	paymentUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type DepositHandler struct {
	invoiceService *billing.InvoiceService
	paymentService *paymentUsecase.PaymentService
}

func NewDepositHandler(invoiceService *billing.InvoiceService, paymentService *paymentUsecase.PaymentService) *DepositHandler {
	return &DepositHandler{
		invoiceService: invoiceService,
		paymentService: paymentService,
	}
}

type depositRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Gateway  string  `json:"gateway"`
}

func (h *DepositHandler) DepositFunds(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	var req depositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Valid deposit amount is required", nil)
		return
	}

	curr := req.Currency
	if curr == "" {
		curr = "USD"
	}

	inv, err := h.invoiceService.CreateInvoice(r.Context(), billing.CreateInvoiceDTO{
		ClientID: clientID,
		Currency: curr,
		DueDays:  7,
		Items: []billing.CreateInvoiceItemDTO{
			{
				Title:    "Account Balance Deposit / Top-up",
				Price:    decimal.FromFloat(req.Amount),
				Quantity: 1,
				Taxable:  false,
			},
		},
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DEPOSIT_FAILED", "Failed to generate deposit invoice", err.Error())
		return
	}

	data := map[string]interface{}{
		"invoice_id": inv.ID,
		"nr":         inv.Nr,
		"total":      inv.Total.ToFloat(),
		"currency":   inv.Currency,
		"status":     inv.Status,
		"message":    "Deposit invoice generated successfully",
	}

	if req.Gateway != "" {
		paymentRes, err := h.paymentService.InitiateInvoicePayment(r.Context(), inv.ID, req.Gateway)
		if err == nil {
			data["payment"] = paymentRes
			data["redirect_url"] = paymentRes.RedirectURL
		}
	}

	response.JSON(w, http.StatusCreated, data, nil)
}
