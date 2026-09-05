package guest

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/currency"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CurrencyHandler struct {
	currencyService *currency.CurrencyService
}

func NewCurrencyHandler(currencyService *currency.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{currencyService: currencyService}
}

func (h *CurrencyHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.currencyService.ListCurrencies(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve currencies", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, list, nil)
}
