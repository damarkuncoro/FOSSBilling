package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	themeUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/theme"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ThemeHandler struct {
	svc *themeUsecase.ThemeService
}

func NewThemeHandler(svc *themeUsecase.ThemeService) *ThemeHandler {
	return &ThemeHandler{svc: svc}
}

func (h *ThemeHandler) ListThemes(w http.ResponseWriter, r *http.Request) {
	target := domain.ThemeTarget(r.URL.Query().Get("target"))
	themes, err := h.svc.ListThemes(r.Context(), target)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve themes", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"list": themes}, nil)
}

func (h *ThemeHandler) GetCurrentTheme(w http.ResponseWriter, r *http.Request) {
	target := domain.ThemeTarget(r.URL.Query().Get("target"))
	if target == "" {
		target = domain.ThemeTargetClient
	}
	theme, err := h.svc.GetCurrentTheme(r.Context(), target)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Current theme not found", nil)
		return
	}
	response.JSON(w, http.StatusOK, theme, nil)
}

func (h *ThemeHandler) SelectTheme(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code   string             `json:"code"`
		Target domain.ThemeTarget `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Invalid JSON payload", err.Error())
		return
	}
	if req.Target == "" {
		req.Target = domain.ThemeTargetClient
	}

	if err := h.svc.SelectTheme(r.Context(), req.Code, req.Target); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Theme code not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to set active theme", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Theme activated successfully",
	}, nil)
}

func (h *ThemeHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	cfg, err := h.svc.GetConfig(r.Context(), code)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Theme config not found", nil)
		return
	}
	response.JSON(w, http.StatusOK, cfg, nil)
}

func (h *ThemeHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Invalid JSON payload", err.Error())
		return
	}

	if err := h.svc.UpdateConfig(r.Context(), code, req); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update theme configuration", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Theme configuration updated successfully",
	}, nil)
}
