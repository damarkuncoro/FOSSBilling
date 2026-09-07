package guest

import (
	"errors"
	"net/http"
	"strconv"

	formbuilderUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type FormbuilderHandler struct {
	svc *formbuilderUsecase.FormbuilderService
}

func NewFormbuilderHandler(svc *formbuilderUsecase.FormbuilderService) *FormbuilderHandler {
	return &FormbuilderHandler{svc: svc}
}

func (h *FormbuilderHandler) GetForm(w http.ResponseWriter, r *http.Request) {
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
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve form", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, form, nil)
}
