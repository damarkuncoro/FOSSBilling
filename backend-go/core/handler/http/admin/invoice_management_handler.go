package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

type InvoiceManagementHandler struct {
	staffService   *staff.StaffService
	invoiceRepo    domain.InvoiceRepository
	clientRepo     domain.ClientRepository
	invoiceService *billing.InvoiceService
}

func NewInvoiceManagementHandler(
	staffService *staff.StaffService,
	invoiceRepo domain.InvoiceRepository,
	clientRepo domain.ClientRepository,
	invoiceService *billing.InvoiceService,
) *InvoiceManagementHandler {
	return &InvoiceManagementHandler{
		staffService:   staffService,
		invoiceRepo:    invoiceRepo,
		clientRepo:     clientRepo,
		invoiceService: invoiceService,
	}
}

func (h *InvoiceManagementHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "invoices", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: invoices", nil)
		return
	}

	limit := 50
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

	invoices, total, err := h.invoiceRepo.List(r.Context(), limit, offset)
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

func (h *InvoiceManagementHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "invoices", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: invoices", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid invoice ID", nil)
		return
	}

	invoice, err := h.invoiceRepo.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Invoice not found", nil)
		return
	}

	response.JSON(w, http.StatusOK, invoice, nil)
}

type createInvoiceItemReq struct {
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Taxable  bool    `json:"taxable"`
}

type createInvoiceReq struct {
	ClientID int64                  `json:"client_id"`
	Currency string                 `json:"currency"`
	DueDays  int                    `json:"due_days"`
	Items    []createInvoiceItemReq `json:"items"`
}

func (h *InvoiceManagementHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "invoices", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: invoices", nil)
		return
	}

	var req createInvoiceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ClientID == 0 || len(req.Items) == 0 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Valid client_id and at least one item are required", nil)
		return
	}

	var itemsDTO []billing.CreateInvoiceItemDTO
	for _, it := range req.Items {
		if it.Quantity <= 0 {
			response.Error(w, http.StatusBadRequest, "VALIDATION_FAILED", "Quantity must be greater than zero for item: "+it.Title, nil)
			return
		}
		itemsDTO = append(itemsDTO, billing.CreateInvoiceItemDTO{
			Title:    security.SanitizeHTML(it.Title),
			Price:    decimal.FromFloat(it.Price),
			Quantity: it.Quantity,
			Taxable:  it.Taxable,
		})
	}

	dueDays := req.DueDays
	if dueDays <= 0 {
		dueDays = 14
	}

	inv, err := h.invoiceService.CreateInvoice(r.Context(), billing.CreateInvoiceDTO{
		ClientID: req.ClientID,
		Currency: req.Currency,
		DueDays:  dueDays,
		Items:    itemsDTO,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVOICE_CREATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, inv, nil)
}

func (h *InvoiceManagementHandler) RefundInvoice(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "invoices", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: invoices", nil)
		return
	}

	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.invoiceService.RefundInvoice(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, "REFUND_FAILED", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}

func (h *InvoiceManagementHandler) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "invoices", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: invoices", nil)
		return
	}

	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.invoiceRepo.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}
