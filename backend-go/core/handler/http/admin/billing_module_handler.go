package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	billing "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type BillingModuleHandler struct {
	staffSvc *staff.StaffService; statsSvc *stats.StatsService; tax *billing.TaxCalculator; promoRepo domain.PromoRepository
}

func NewBillingModuleHandler(s *staff.StaffService, st *stats.StatsService, t *billing.TaxCalculator, p domain.PromoRepository) *BillingModuleHandler {
	return &BillingModuleHandler{s, st, t, p}
}

func (h *BillingModuleHandler) ListGateways(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "read") { return }
	response.JSON(w, 200, []any{map[string]any{"id": "stripe", "enabled": true}}, nil)
}

func (h *BillingModuleHandler) ListTaxRules(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "read") { return }
	rs, err := h.tax.ListRules(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, rs, nil)
}

func (h *BillingModuleHandler) CreateTaxRule(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "write") { return }
	var rule domain.TaxRule; if request.Decode(r, &rule) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.tax.CreateRule(r.Context(), &rule); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, rule, nil)
}

func (h *BillingModuleHandler) UpdateTaxRule(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "write") { return }
	var rule domain.TaxRule; if request.Decode(r, &rule) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	rule.ID = request.GetID(r)
	if err := h.tax.UpdateRule(r.Context(), &rule); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, rule, nil)
}

func (h *BillingModuleHandler) DeleteTaxRule(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "delete") { return }
	_ = h.tax.DeleteRule(r.Context(), request.GetID(r))
	response.JSON(w, 200, map[string]bool{"deleted": true}, nil)
}

func (h *BillingModuleHandler) ListCoupons(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "read") { return }
	cs, _, err := h.promoRepo.List(r.Context(), 100, 0)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, cs, nil)
}

func (h *BillingModuleHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "write") { return }
	var p domain.Promo; if request.Decode(r, &p) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.promoRepo.Create(r.Context(), &p); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, p, nil)
}

func (h *BillingModuleHandler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "delete") { return }
	_ = h.promoRepo.Delete(r.Context(), request.GetID(r))
	response.JSON(w, 200, map[string]bool{"deleted": true}, nil)
}

func (h *BillingModuleHandler) ListEmailTemplates(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, []any{}, nil) }
func (h *BillingModuleHandler) GetMailConfig(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, map[string]any{}, nil) }
func (h *BillingModuleHandler) SendTestEmail(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, map[string]any{"ok": true}, nil) }

func (h *BillingModuleHandler) GetFinancialReports(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "billing", "read") { return }
	rep, err := h.statsSvc.GetFinancialReports(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, rep, nil)
}

func (h *BillingModuleHandler) ExportInvoicesCSV(w http.ResponseWriter, r *http.Request) {
	csv, err := h.statsSvc.GenerateInvoicesCSV(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	w.Header().Set("Content-Type", "text/csv"); w.Write([]byte(csv))
}
