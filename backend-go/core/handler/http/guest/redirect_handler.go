package guest

import (
	"errors"
	"net/http"

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

func (h *RedirectHandler) Lookup(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_PATH", "Path query parameter is required", nil)
		return
	}

	rule, err := h.svc.MatchRedirect(r.Context(), path)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "No active redirect found for path", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to lookup redirect", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"path":        rule.Path,
		"target":      rule.Target,
		"status_code": rule.StatusCode,
	}, nil)
}
