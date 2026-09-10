package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/currency"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CurrencyHandler struct {
	svc *currency.CurrencyService; staffSvc *staff.StaffService
}

func NewCurrencyHandler(s *currency.CurrencyService, ss *staff.StaffService) *CurrencyHandler { return &CurrencyHandler{s, ss} }

func (h *CurrencyHandler) List(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "read") { return }
	l, err := h.svc.ListCurrencies(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, l, nil)
}

func (h *CurrencyHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "write") { return }
	var d currency.CreateCurrencyDTO; if request.Decode(r, &d) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	c, err := h.svc.CreateCurrency(r.Context(), d)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, c, nil)
}

func (h *CurrencyHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "write") { return }
	var d currency.UpdateCurrencyDTO; if request.Decode(r, &d) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	c, err := h.svc.UpdateCurrency(r.Context(), request.GetParam(r, "code"), d)
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, c, nil)
}

func (h *CurrencyHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "write") { return }
	if err := h.svc.SetDefault(r.Context(), request.GetParam(r, "code")); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"status": "success"}, nil)
}

func (h *CurrencyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "write") { return }
	if err := h.svc.DeleteCurrency(r.Context(), request.GetParam(r, "code")); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"status": "deleted"}, nil)
}

func (h *CurrencyHandler) SyncRates(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "system", "write") { return }
	if err := h.svc.UpdateExchangeRates(r.Context()); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"status": "success"}, nil)
}
