package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	redirectUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/redirect"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type RedirectHandler struct {
	svc *redirectUsecase.RedirectService
}

func NewRedirectHandler(svc *redirectUsecase.RedirectService) *RedirectHandler {
	return &RedirectHandler{svc: svc}
}

func (h *RedirectHandler) ListRedirects(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	list, total, err := h.svc.ListRedirects(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve redirects", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"list":  list,
		"total": total,
	}, nil)
}

func (h *RedirectHandler) GetRedirect(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid redirect ID", nil)
		return
	}

	item, err := h.svc.GetRedirect(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Redirect rule not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get redirect", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, item, nil)
}

func (h *RedirectHandler) CreateRedirect(w http.ResponseWriter, r *http.Request) {
	var req redirectUsecase.CreateRedirectDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Invalid JSON payload", err.Error())
		return
	}

	created, err := h.svc.CreateRedirect(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, created, nil)
}

func (h *RedirectHandler) UpdateRedirect(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid redirect ID", nil)
		return
	}

	var req redirectUsecase.UpdateRedirectDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Invalid JSON payload", err.Error())
		return
	}
	req.ID = id

	updated, err := h.svc.UpdateRedirect(r.Context(), req)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Redirect rule not found", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, updated, nil)
}

func (h *RedirectHandler) DeleteRedirect(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid redirect ID", nil)
		return
	}

	if err := h.svc.DeleteRedirect(r.Context(), id); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Redirect rule not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete redirect", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Redirect rule deleted successfully",
	}, nil)
}
