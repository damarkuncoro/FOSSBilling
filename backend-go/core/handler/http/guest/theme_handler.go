package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	themeUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/theme"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type ThemeHandler struct {
	svc *themeUsecase.ThemeService
}

func NewThemeHandler(svc *themeUsecase.ThemeService) *ThemeHandler {
	return &ThemeHandler{svc: svc}
}

func (h *ThemeHandler) GetActiveTheme(w http.ResponseWriter, r *http.Request) {
	theme, err := h.svc.GetCurrentTheme(r.Context(), domain.ThemeTargetClient)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Active client theme not found", nil)
		return
	}
	response.JSON(w, http.StatusOK, theme, nil)
}
