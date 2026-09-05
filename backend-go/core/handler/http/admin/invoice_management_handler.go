package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type InvoiceManagementHandler struct {
	invoiceRepo    domain.InvoiceRepository
	clientRepo     domain.ClientRepository
	invoiceService *billing.InvoiceService
}

func NewInvoiceManagementHandler(
	invoiceRepo domain.InvoiceRepository,
	clientRepo domain.ClientRepository,
	invoiceService *billing.InvoiceService,
) *InvoiceManagementHandler {
	return &InvoiceManagementHandler{
		invoiceRepo:    invoiceRepo,
		clientRepo:     clientRepo,
		invoiceService: invoiceService,
	}
}

func (h *InvoiceManagementHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
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
	role := middleware.GetRole(r.Context())
	if role != "admin" && role != "superadmin" && role != "billing" && role != "staff" {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Staff permission required", nil)
		return
	}

	var req createInvoiceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ClientID == 0 || len(req.Items) == 0 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Valid client_id and at least one item are required", nil)
		return
	}

	var itemsDTO []billing.CreateInvoiceItemDTO
	for _, it := range req.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		itemsDTO = append(itemsDTO, billing.CreateInvoiceItemDTO{
			Title:    it.Title,
			Price:    decimal.FromFloat(it.Price),
			Quantity: qty,
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
