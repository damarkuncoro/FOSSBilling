package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	formbuilderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type FormbuilderHandler struct {
	staffService *staff.StaffService
	svc          *formbuilderUsecase.FormbuilderService
}

func NewFormbuilderHandler(staffService *staff.StaffService, svc *formbuilderUsecase.FormbuilderService) *FormbuilderHandler {
	return &FormbuilderHandler{staffService: staffService, svc: svc}
}

func (h *FormbuilderHandler) ListForms(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	forms, total, err := h.svc.ListForms(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve forms", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"list":  forms,
		"total": total,
	}, nil)
}

func (h *FormbuilderHandler) GetForm(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid form ID", nil)
		return
	}

	form, err := h.svc.GetForm(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Form not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get form", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, form, nil)
}

func (h *FormbuilderHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	var req formbuilderUsecase.CreateFormDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	form, err := h.svc.CreateForm(r.Context(), req)
	if err != nil {
		if errors.Is(err, appErrors.ErrInvalidInput) {
			response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create form", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, form, nil)
}

func (h *FormbuilderHandler) UpdateForm(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid form ID", nil)
		return
	}

	var req formbuilderUsecase.UpdateFormDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	form, err := h.svc.UpdateForm(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Form not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update form", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, form, nil)
}

func (h *FormbuilderHandler) DeleteForm(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid form ID", nil)
		return
	}

	if err := h.svc.DeleteForm(r.Context(), id); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Form not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete form", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Form successfully deleted"}, nil)
}

func (h *FormbuilderHandler) AddField(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	formIDStr := r.PathValue("id")
	formID, err := strconv.ParseInt(formIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid form ID", nil)
		return
	}

	var req formbuilderUsecase.FormFieldDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	field, err := h.svc.AddField(r.Context(), formID, req)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Form not found", nil)
			return
		}
		if errors.Is(err, formbuilderUsecase.ErrInvalidFieldType) ||
			errors.Is(err, formbuilderUsecase.ErrFieldNameNumeric) ||
			errors.Is(err, formbuilderUsecase.ErrFieldNameEmpty) ||
			errors.Is(err, formbuilderUsecase.ErrDuplicateFieldName) {
			response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add form field", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, field, nil)
}

func (h *FormbuilderHandler) UpdateField(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	fieldIDStr := r.PathValue("field_id")
	fieldID, err := strconv.ParseInt(fieldIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid field ID", nil)
		return
	}

	var req formbuilderUsecase.FormFieldDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	field, err := h.svc.UpdateField(r.Context(), fieldID, req)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Form field not found", nil)
			return
		}
		if errors.Is(err, formbuilderUsecase.ErrInvalidFieldType) ||
			errors.Is(err, formbuilderUsecase.ErrFieldNameNumeric) ||
			errors.Is(err, formbuilderUsecase.ErrFieldNameEmpty) ||
			errors.Is(err, formbuilderUsecase.ErrDuplicateFieldName) {
			response.Error(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", err.Error(), nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update form field", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, field, nil)
}

func (h *FormbuilderHandler) DeleteField(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "formbuilder", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: formbuilder", nil)
		return
	}

	fieldIDStr := r.PathValue("field_id")
	fieldID, err := strconv.ParseInt(fieldIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid field ID", nil)
		return
	}

	if err := h.svc.DeleteField(r.Context(), fieldID); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Form field not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete form field", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Form field successfully deleted"}, nil)
}
