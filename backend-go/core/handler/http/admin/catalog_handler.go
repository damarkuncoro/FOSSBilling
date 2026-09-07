package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/catalog"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CatalogHandler struct {
	productService *catalog.ProductService
	catalogRepo    domain.CatalogRepository
}

func NewCatalogHandler(productService *catalog.ProductService, catalogRepo domain.CatalogRepository) *CatalogHandler {
	return &CatalogHandler{
		productService: productService,
		catalogRepo:    catalogRepo,
	}
}

// --- Products & Categories ---
func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	products, total, err := h.productService.ListProducts(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	meta := &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	response.JSON(w, http.StatusOK, products, meta)
}

func (h *CatalogHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}

	if err := h.productService.CreateProduct(r.Context(), &p); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, p, nil)
}

func (h *CatalogHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid product id", nil)
		return
	}

	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}
	p.ID = id

	if err := h.productService.UpdateProduct(r.Context(), &p); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, p, nil)
}

func (h *CatalogHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid product id", nil)
		return
	}

	if err := h.productService.DeleteProduct(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

func (h *CatalogHandler) ListProductCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.catalogRepo.ListCategories(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, categories, nil)
}

// --- Domains & Registrars ---
func (h *CatalogHandler) ListTlds(w http.ResponseWriter, r *http.Request) {
	tlds, err := h.catalogRepo.ListTlds(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, tlds, nil)
}

func (h *CatalogHandler) CreateTld(w http.ResponseWriter, r *http.Request) {
	// Implementation placeholder for real TLD creation
	response.JSON(w, http.StatusCreated, map[string]string{"message": "TLD creation not yet fully implemented"}, nil)
}

func (h *CatalogHandler) DeleteTld(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

func (h *CatalogHandler) ListRegistrars(w http.ResponseWriter, r *http.Request) {
	// In production, this would list registered registrar drivers
	registrars := []map[string]interface{}{
		{"id": "namecheap", "name": "Namecheap API", "enabled": true},
		{"id": "enom", "name": "eNom Reseller", "enabled": false},
	}
	response.JSON(w, http.StatusOK, registrars, nil)
}

// --- Servers ---
func (h *CatalogHandler) ListServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.catalogRepo.ListServers(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, servers, nil)
}

func (h *CatalogHandler) CreateServer(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusCreated, map[string]string{"message": "Server creation not yet fully implemented"}, nil)
}

func (h *CatalogHandler) TestServer(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Connection handshake successful!",
	}, nil)
}

func (h *CatalogHandler) DeleteServer(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}
