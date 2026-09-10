package client

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	payment "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type InvoiceHandler struct {
	ir domain.InvoiceRepository; cr domain.ClientRepository; cor domain.CompanyRepository; is *billing.InvoiceService; ps *payment.PaymentService
}

func NewInvoiceHandler(ir domain.InvoiceRepository, cr domain.ClientRepository, cor domain.CompanyRepository, is *billing.InvoiceService, ps *payment.PaymentService) *InvoiceHandler {
	return &InvoiceHandler{ir, cr, cor, is, ps}
}

func (h *InvoiceHandler) check(w http.ResponseWriter, r *http.Request) (int64, *domain.Invoice, bool) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return 0, nil, false }
	i, err := h.ir.GetByID(r.Context(), request.GetID(r))
	if err != nil || i.ClientID != cid { response.Error(w, 403, "FORBIDDEN", "No access", nil); return 0, nil, false }
	return cid, i, true
}

func (h *InvoiceHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	l, o := request.GetLimitOffset(r)
	is, tot, err := h.ir.ListByClientID(r.Context(), cid, l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, is, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	_, i, ok := h.check(w, r); if ok { response.JSON(w, 200, i, nil) }
}

func (h *InvoiceHandler) PayWithBalance(w http.ResponseWriter, r *http.Request) {
	cid, i, ok := h.check(w, r); if !ok { return }
	upd, err := h.is.PayWithBalance(r.Context(), cid, i.ID)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, upd, nil)
}

func (h *InvoiceHandler) PayWithGateway(w http.ResponseWriter, r *http.Request) {
	_, i, ok := h.check(w, r); if !ok { return }
	var req struct{ Gateway string `json:"gateway"` }; _ = request.Decode(r, &req)
	res, err := h.ps.InitiateInvoicePayment(r.Context(), i.ID, req.Gateway)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, res, nil)
}

func (h *InvoiceHandler) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); adm := middleware.GetRole(r.Context()) != "client"
	invID := request.GetID(r)

	// Access check
	i, err := h.ir.GetByID(r.Context(), invID)
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Invoice not found", nil); return }
	if cid != i.ClientID && !adm { response.Error(w, 403, "FORBIDDEN", "No access", nil); return }

	data, filename, err := h.is.GeneratePDF(r.Context(), invID)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	_, _ = w.Write(data)
}
