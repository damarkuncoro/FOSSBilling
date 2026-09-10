package admin

import (
	"fmt"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/catalog"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CatalogHandler struct {
	staffService *staff.StaffService; productService *catalog.ProductService; serverService *catalog.ServerService; catalogRepo domain.CatalogRepository
}

func NewCatalogHandler(s *staff.StaffService, ps *catalog.ProductService, ss *catalog.ServerService, cr domain.CatalogRepository) *CatalogHandler {
	return &CatalogHandler{s, ps, ss, cr}
}

func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "products", "read") { return }
	l, o := request.GetLimitOffset(r)
	ps, tot, err := h.productService.ListProducts(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ps, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *CatalogHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "products", "write") { return }
	var p domain.Product; if request.Decode(r, &p) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.productService.CreateProduct(r.Context(), &p); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, p, nil)
}

func (h *CatalogHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "products", "write") { return }
	var p domain.Product; if request.Decode(r, &p) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	p.ID = request.GetID(r)
	if err := h.productService.UpdateProduct(r.Context(), &p); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, p, nil)
}

func (h *CatalogHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "products", "delete") { return }
	if err := h.productService.DeleteProduct(r.Context(), request.GetID(r)); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"deleted": true}, nil)
}

func (h *CatalogHandler) ListProductCategories(w http.ResponseWriter, r *http.Request) {
	cs, err := h.catalogRepo.ListCategories(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, cs, nil)
}

func (h *CatalogHandler) ListTlds(w http.ResponseWriter, r *http.Request) {
	ts, err := h.catalogRepo.ListTlds(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ts, nil)
}

func (h *CatalogHandler) CreateTld(w http.ResponseWriter, r *http.Request) { response.JSON(w, 201, map[string]string{"m": "OK"}, nil) }
func (h *CatalogHandler) DeleteTld(w http.ResponseWriter, r *http.Request) { response.JSON(w, 200, map[string]bool{"deleted": true}, nil) }
func (h *CatalogHandler) ListRegistrars(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, 200, []any{map[string]any{"id": "namecheap", "name": "Namecheap"}}, nil)
}

func (h *CatalogHandler) ListServers(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "servers", "read") { return }
	ss, err := h.serverService.ListServers(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ss, nil)
}

func (h *CatalogHandler) CreateServer(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "servers", "write") { return }
	var s domain.Server; if request.Decode(r, &s) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.serverService.CreateServer(r.Context(), &s); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, s, nil)
}

func (h *CatalogHandler) TestServer(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "servers", "write") { return }
	err := h.serverService.TestConnection(r.Context(), request.GetID(r))
	response.JSON(w, 200, map[string]any{"success": err == nil, "message": fmt.Sprintf("%v", err)}, nil)
}

func (h *CatalogHandler) DeleteServer(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "servers", "delete") { return }
	if err := h.serverService.DeleteServer(r.Context(), request.GetID(r)); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"deleted": true}, nil)
}
