package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	extension "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/extension"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ExtensionHandler struct {
	staffService *staff.StaffService; svc *extension.ExtensionService
}

func NewExtensionHandler(ss *staff.StaffService, s *extension.ExtensionService) *ExtensionHandler { return &ExtensionHandler{ss, s} }

func (h *ExtensionHandler) ListExtensions(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "read") { return }
	f := domain.ExtensionFilter{Type: r.URL.Query().Get("type"), Status: r.URL.Query().Get("status"), Search: r.URL.Query().Get("search")}
	if r.URL.Query().Has("active") { a := r.URL.Query().Get("active") == "true"; f.Active = &a }
	l, err := h.svc.ListExtensions(r.Context(), f)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, l, nil)
}

func (h *ExtensionHandler) GetExtension(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "read") { return }
	e, err := h.svc.GetExtension(r.Context(), request.GetParam(r, "id"))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, e, nil)
}

func (h *ExtensionHandler) Activate(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "write") { return }
	e, err := h.svc.ActivateExtension(r.Context(), request.GetParam(r, "id"))
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, e, nil)
}

func (h *ExtensionHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "write") { return }
	e, err := h.svc.DeactivateExtension(r.Context(), request.GetParam(r, "id"))
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, e, nil)
}

func (h *ExtensionHandler) Install(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "write") { return }
	e, err := h.svc.InstallExtension(r.Context(), request.GetParam(r, "id"))
	if err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, e, nil)
}

func (h *ExtensionHandler) Uninstall(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "delete") { return }
	if err := h.svc.UninstallExtension(r.Context(), request.GetParam(r, "id")); err != nil { response.Error(w, 400, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]any{"ok": true}, nil)
}

func (h *ExtensionHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "read") { return }
	c, err := h.svc.GetExtensionConfig(r.Context(), request.GetParam(r, "id"))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, c, nil)
}

func (h *ExtensionHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "extensions", "write") { return }
	var c map[string]any; if request.Decode(r, &c) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.svc.UpdateExtensionConfig(r.Context(), request.GetParam(r, "id"), c); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]any{"ok": true}, nil)
}

func (h *ExtensionHandler) ListMarketplace(w http.ResponseWriter, r *http.Request) {
	l, err := h.svc.FetchMarketplace(r.Context(), r.URL.Query().Get("type"))
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, l, nil)
}

func (h *ExtensionHandler) GetMarketplaceReadme(w http.ResponseWriter, r *http.Request) {
	rm, err := h.svc.GetMarketplaceReadme(r.Context(), request.GetParam(r, "id"))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, map[string]string{"readme": rm}, nil)
}
