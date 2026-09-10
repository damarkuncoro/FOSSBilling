package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type EmailTemplateHandler struct {
	staffService *staff.StaffService
	repo         domain.EmailTemplateRepository
}

func NewEmailTemplateHandler(s *staff.StaffService, r domain.EmailTemplateRepository) *EmailTemplateHandler {
	return &EmailTemplateHandler{s, r}
}

func (h *EmailTemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "read") {
		return
	}
	ts, err := h.repo.List(r.Context())
	if err != nil {
		response.Error(w, 500, "ERR", err.Error(), nil)
		return
	}
	response.JSON(w, 200, ts, nil)
}

func (h *EmailTemplateHandler) Get(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "read") {
		return
	}
	t, err := h.repo.GetByCode(r.Context(), r.PathValue("code"))
	if err != nil {
		response.Error(w, 404, "NOT_FOUND", "Template not found", nil)
		return
	}
	response.JSON(w, 200, t, nil)
}

func (h *EmailTemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "system", "write") {
		return
	}
	var t domain.EmailTemplate
	if err := request.Decode(r, &t); err != nil {
		response.Error(w, 400, "BAD", "Invalid request body", nil)
		return
	}

	if err := h.repo.Update(r.Context(), &t); err != nil {
		response.Error(w, 500, "ERR", err.Error(), nil)
		return
	}
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}
