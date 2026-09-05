package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/currency"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/staff"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CurrencyHandler struct {
	currencyService *currency.CurrencyService
	staffService    *staff.StaffService
}

func NewCurrencyHandler(currencyService *currency.CurrencyService, staffService *staff.StaffService) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
		staffService:    staffService,
	}
}

func (h *CurrencyHandler) List(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	list, err := h.currencyService.ListCurrencies(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve currencies", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, list, nil)
}

func (h *CurrencyHandler) Create(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	var dto currency.CreateCurrencyDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	c, err := h.currencyService.CreateCurrency(r.Context(), dto)
	if err != nil {
		if errors.Is(err, appErrors.ErrDuplicateEntry) {
			response.Error(w, http.StatusConflict, "DUPLICATE_CURRENCY", "Currency code already exists", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, c, nil)
}

func (h *CurrencyHandler) Update(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	code := parts[len(parts)-1]

	var dto currency.UpdateCurrencyDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	c, err := h.currencyService.UpdateCurrency(r.Context(), code, dto)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Currency not found", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, c, nil)
}

func (h *CurrencyHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	code := parts[len(parts)-2]

	if err := h.currencyService.SetDefault(r.Context(), code); err != nil {
		response.Error(w, http.StatusBadRequest, "ACTION_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Default base currency updated to " + strings.ToUpper(code),
	}, nil)
}

func (h *CurrencyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "system", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: system", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	code := parts[len(parts)-1]

	if err := h.currencyService.DeleteCurrency(r.Context(), code); err != nil {
		response.Error(w, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "deleted",
		"message": "Currency deleted successfully",
	}, nil)
}
