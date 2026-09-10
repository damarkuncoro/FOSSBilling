package admin

import (
	"net/http"

	formbuilder "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type FormbuilderHandler struct {
	staffSvc *staff.StaffService; svc *formbuilder.FormbuilderService
}

func NewFormbuilderHandler(ss *staff.StaffService, s *formbuilder.FormbuilderService) *FormbuilderHandler { return &FormbuilderHandler{ss, s} }

func (h *FormbuilderHandler) ListForms(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "read") { return }
	l, o := request.GetLimitOffset(r)
	fs, tot, err := h.svc.ListForms(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, fs, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *FormbuilderHandler) GetForm(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "read") { return }
	f, err := h.svc.GetForm(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, f, nil)
}

func (h *FormbuilderHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "write") { return }
	var req formbuilder.CreateFormDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	f, err := h.svc.CreateForm(r.Context(), req)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, f, nil)
}

func (h *FormbuilderHandler) UpdateForm(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "write") { return }
	var req formbuilder.UpdateFormDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	f, err := h.svc.UpdateForm(r.Context(), request.GetID(r), req)
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, f, nil)
}

func (h *FormbuilderHandler) DeleteForm(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "delete") { return }
	if err := h.svc.DeleteForm(r.Context(), request.GetID(r)); err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, map[string]string{"message": "Deleted"}, nil)
}

func (h *FormbuilderHandler) AddField(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "write") { return }
	var req formbuilder.FormFieldDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	f, err := h.svc.AddField(r.Context(), request.GetID(r), req)
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, f, nil)
}

func (h *FormbuilderHandler) UpdateField(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "write") { return }
	var req formbuilder.FormFieldDTO; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	f, err := h.svc.UpdateField(r.Context(), request.GetID(r), req)
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, f, nil)
}

func (h *FormbuilderHandler) DeleteField(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffSvc, "formbuilder", "delete") { return }
	if err := h.svc.DeleteField(r.Context(), request.GetID(r)); err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, map[string]string{"message": "Deleted"}, nil)
}
