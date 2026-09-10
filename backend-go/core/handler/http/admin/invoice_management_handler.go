package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type InvoiceManagementHandler struct {
	staffService *staff.StaffService; invoiceRepo domain.InvoiceRepository; clientRepo domain.ClientRepository; invoiceService *billing.InvoiceService
}

func NewInvoiceManagementHandler(s *staff.StaffService, ir domain.InvoiceRepository, cr domain.ClientRepository, is *billing.InvoiceService) *InvoiceManagementHandler {
	return &InvoiceManagementHandler{s, ir, cr, is}
}

func (h *InvoiceManagementHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "invoices", "read") { return }
	l, o := request.GetLimitOffset(r)
	invs, tot, err := h.invoiceRepo.List(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, invs, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *InvoiceManagementHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "invoices", "read") { return }
	inv, err := h.invoiceRepo.GetByID(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, inv, nil)
}

type itReq struct {
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Taxable  bool    `json:"taxable"`
}

type invReq struct {
	ClientID int64   `json:"client_id"`
	Currency string  `json:"currency"`
	DueDays  int     `json:"due_days"`
	Items    []itReq `json:"items"`
}

func (h *InvoiceManagementHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "invoices", "write") { return }
	var req invReq; if request.Decode(r, &req) != nil || req.ClientID == 0 || len(req.Items) == 0 { response.Error(w, 400, "BAD", "Invalid", nil); return }
	var its []billing.CreateInvoiceItemDTO
	for _, i := range req.Items { its = append(its, billing.CreateInvoiceItemDTO{Title: i.Title, Price: decimal.FromFloat(i.Price), Quantity: i.Quantity, Taxable: i.Taxable}) }
	inv, err := h.invoiceService.CreateInvoice(r.Context(), billing.CreateInvoiceDTO{ClientID: req.ClientID, Currency: req.Currency, DueDays: req.DueDays, Items: its})
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, inv, nil)
}

func (h *InvoiceManagementHandler) RefundInvoice(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "invoices", "write") { return }
	if err := h.invoiceService.RefundInvoice(r.Context(), request.GetID(r)); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *InvoiceManagementHandler) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "invoices", "delete") { return }
	if err := h.invoiceRepo.Delete(r.Context(), request.GetID(r)); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *InvoiceManagementHandler) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "invoices", "read") { return }
	data, filename, err := h.invoiceService.GeneratePDF(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	_, _ = w.Write(data)
}
