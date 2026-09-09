package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	extensionUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/extension"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ExtensionHandler struct {
	staffService *staff.StaffService
	svc          *extensionUsecase.ExtensionService
}

func NewExtensionHandler(staffService *staff.StaffService, svc *extensionUsecase.ExtensionService) *ExtensionHandler {
	return &ExtensionHandler{staffService: staffService, svc: svc}
}

func (h *ExtensionHandler) ListExtensions(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	var filter domain.ExtensionFilter
	filter.Type = r.URL.Query().Get("type")
	filter.Status = r.URL.Query().Get("status")
	filter.Search = r.URL.Query().Get("search")

	if r.URL.Query().Has("active") {
		active := r.URL.Query().Get("active") == "true" || r.URL.Query().Get("active") == "1"
		filter.Active = &active
	}
	if r.URL.Query().Has("has_settings") {
		hasSettings := r.URL.Query().Get("has_settings") == "true" || r.URL.Query().Get("has_settings") == "1"
		filter.HasSettings = &hasSettings
	}

	list, err := h.svc.ListExtensions(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list extensions", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, list, nil)
}

func (h *ExtensionHandler) GetExtension(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	id := r.PathValue("id")
	ext, err := h.svc.GetExtension(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Extension not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve extension", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, ext, nil)
}

func (h *ExtensionHandler) Activate(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	id := r.PathValue("id")
	ext, err := h.svc.ActivateExtension(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Extension not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to activate extension", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, ext, nil)
}

func (h *ExtensionHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	id := r.PathValue("id")
	ext, err := h.svc.DeactivateExtension(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Extension not found", nil)
			return
		}
		if errors.Is(err, extensionUsecase.ErrCoreExtensionImmutable) {
			response.Error(w, http.StatusForbidden, "IMMUTABLE_EXTENSION", err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to deactivate extension", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, ext, nil)
}

func (h *ExtensionHandler) Install(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	id := r.PathValue("id")
	ext, err := h.svc.InstallExtension(r.Context(), id)
	if err != nil {
		if errors.Is(err, extensionUsecase.ErrExtensionNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Marketplace extension not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to install extension", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, ext, nil)
}

func (h *ExtensionHandler) Uninstall(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	id := r.PathValue("id")
	err := h.svc.UninstallExtension(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Extension not found", nil)
			return
		}
		if errors.Is(err, extensionUsecase.ErrCoreExtensionImmutable) {
			response.Error(w, http.StatusForbidden, "IMMUTABLE_EXTENSION", err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to uninstall extension", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Extension successfully uninstalled"}, nil)
}

func (h *ExtensionHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	id := r.PathValue("id")
	cfg, err := h.svc.GetExtensionConfig(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Extension not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get extension configuration", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, cfg, nil)
}

func (h *ExtensionHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "extensions", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: extensions", nil)
		return
	}

	id := r.PathValue("id")
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	err := h.svc.UpdateExtensionConfig(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Extension not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update extension configuration", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Configuration successfully saved"}, nil)
}

func (h *ExtensionHandler) ListMarketplace(w http.ResponseWriter, r *http.Request) {
	extType := r.URL.Query().Get("type")
	list, err := h.svc.FetchMarketplace(r.Context(), extType)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch marketplace", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, list, nil)
}

func (h *ExtensionHandler) GetMarketplaceReadme(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	readme, err := h.svc.GetMarketplaceReadme(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Marketplace item not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch readme", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"readme": readme}, nil)
}
