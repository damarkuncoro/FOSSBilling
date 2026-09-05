package client

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/pdf"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)


type InvoiceHandler struct {
	invoiceRepo    domain.InvoiceRepository
	clientRepo     domain.ClientRepository
	invoiceService *billing.InvoiceService
}

func NewInvoiceHandler(invoiceRepo domain.InvoiceRepository, clientRepo domain.ClientRepository, invoiceService *billing.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceRepo:    invoiceRepo,
		clientRepo:     clientRepo,
		invoiceService: invoiceService,
	}
}


func (h *InvoiceHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	invoices, total, err := h.invoiceRepo.ListByClientID(r.Context(), clientID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve invoices", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, invoices, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	invoiceIDStr := parts[len(parts)-1]
	invoiceID, err := strconv.ParseInt(invoiceIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid invoice ID", nil)
		return
	}

	invoice, err := h.invoiceRepo.GetByID(r.Context(), invoiceID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Invoice not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get invoice", err.Error())
		return
	}

	if invoice.ClientID != clientID {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not have access to this invoice", nil)
		return
	}

	response.JSON(w, http.StatusOK, invoice, nil)
}

func (h *InvoiceHandler) PayWithBalance(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	// Path e.g. /api/v1/client/invoices/123/pay-balance
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	invoiceIDStr := parts[len(parts)-2]
	invoiceID, err := strconv.ParseInt(invoiceIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid invoice ID", nil)
		return
	}

	invoice, err := h.invoiceRepo.GetByID(r.Context(), invoiceID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Invoice not found", nil)
		return
	}
	if invoice.ClientID != clientID {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not have access to this invoice", nil)
		return
	}

	updated, err := h.invoiceService.PayWithBalance(r.Context(), invoiceID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "PAYMENT_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":  "paid",
		"message": "Invoice paid with account balance successfully",
		"invoice": updated,
	}, nil)
}

func (h *InvoiceHandler) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	role := middleware.GetRole(r.Context())
	isAdmin := role == "admin" || role == "superadmin" || role == "staff" || role == "billing" || role == "support"

	if clientID == 0 && !isAdmin {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	invoiceIDStr := parts[len(parts)-2]
	invoiceID, err := strconv.ParseInt(invoiceIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid invoice ID", nil)
		return
	}

	invoice, err := h.invoiceRepo.GetByID(r.Context(), invoiceID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Invoice not found", nil)
		return
	}
	if !isAdmin && invoice.ClientID != clientID {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not have access to this invoice", nil)
		return
	}

	client, err := h.clientRepo.GetByID(r.Context(), invoice.ClientID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Client not found", nil)
		return
	}

	pdfContent, err := pdf.GenerateInvoicePDF(invoice, client, "", "", "")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "RENDER_ERROR", "Failed to generate invoice PDF", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Invoice-%s.pdf", invoice.Nr))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfContent)
}

