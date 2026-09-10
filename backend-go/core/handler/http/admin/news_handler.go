package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/news"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type NewsHandler struct {
	svc *news.NewsService; staffSvc *staff.StaffService
}

func NewNewsHandler(s *news.NewsService, ss *staff.StaffService) *NewsHandler { return &NewsHandler{s, ss} }

func (h *NewsHandler) List(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "news", "read") { return }
	l, o := request.GetLimitOffset(r)
	ps, tot, err := h.svc.ListAll(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ps, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *NewsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "news", "write") { return }
	var d news.CreateNewsDTO; if request.Decode(r, &d) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	d.AdminID = middleware.GetClientID(r.Context())
	p, err := h.svc.Create(r.Context(), d)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, p, nil)
}

func (h *NewsHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "news", "write") { return }
	var d news.UpdateNewsDTO; if request.Decode(r, &d) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	p, err := h.svc.Update(r.Context(), request.GetID(r), d)
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, p, nil)
}

func (h *NewsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "news", "write") { return }
	if err := h.svc.Delete(r.Context(), request.GetID(r)); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"status": "deleted"}, nil)
}
