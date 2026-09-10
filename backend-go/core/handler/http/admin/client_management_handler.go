package admin

import (
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	pkgAuth "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

type ClientManagementHandler struct {
	staffService *staff.StaffService; authService *auth.AuthUsecase; clientRepo domain.ClientRepository
}

func NewClientManagementHandler(s *staff.StaffService, a *auth.AuthUsecase, cr domain.ClientRepository) *ClientManagementHandler {
	return &ClientManagementHandler{s, a, cr}
}

func (h *ClientManagementHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "clients", "read") { return }
	l, o := request.GetLimitOffset(r)
	cls, tot, err := h.clientRepo.List(r.Context(), l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, cls, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *ClientManagementHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "clients", "read") { return }
	c, err := h.clientRepo.GetByID(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	response.JSON(w, 200, c, nil)
}

type ClientPayload struct { FirstName, LastName, Email, Password, Company, Country, Currency, Status string }

func (h *ClientManagementHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "clients", "write") { return }
	var req ClientPayload; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if req.FirstName == "" || req.Email == "" { response.Error(w, 400, "BAD", "Missing fields", nil); return }
	ph, _ := pkgAuth.HashPassword(req.Password); if req.Password == "" { ph, _ = pkgAuth.HashPassword("Password123!") }
	c := &domain.Client{Email: req.Email, PasswordHash: ph, FirstName: security.SanitizeAlphaNumeric(req.FirstName), LastName: security.SanitizeAlphaNumeric(req.LastName), Company: security.SanitizeHTML(req.Company), Country: security.SanitizeAlphaNumeric(req.Country), Currency: req.Currency, Status: domain.ClientStatus(req.Status), CreatedAt: time.Now().UTC()}
	if c.Currency == "" { c.Currency = "USD" }; if c.Status == "" { c.Status = "active" }
	if err := h.clientRepo.Create(r.Context(), c); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, c, nil)
}

func (h *ClientManagementHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "clients", "write") { return }
	c, err := h.clientRepo.GetByID(r.Context(), request.GetID(r)); if err != nil { response.Error(w, 404, "NOT_FOUND", "Not found", nil); return }
	var req ClientPayload; if request.Decode(r, &req) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if req.FirstName != "" { c.FirstName = security.SanitizeAlphaNumeric(req.FirstName) }
	if req.LastName != "" { c.LastName = security.SanitizeAlphaNumeric(req.LastName) }
	if req.Email != "" { c.Email = req.Email }; if req.Company != "" { c.Company = req.Company }
	if req.Country != "" { c.Country = req.Country }; if req.Currency != "" { c.Currency = req.Currency }
	if req.Status != "" { c.Status = domain.ClientStatus(req.Status) }
	if req.Password != "" { c.PasswordHash, _ = pkgAuth.HashPassword(req.Password) }
	if err := h.clientRepo.Update(r.Context(), c); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, c, nil)
}

func (h *ClientManagementHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "clients", "delete") { return }
	if err := h.clientRepo.Delete(r.Context(), request.GetID(r)); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"deleted": true}, nil)
}

func (h *ClientManagementHandler) ImpersonateClient(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "clients", "write") { return }
	t, err := h.authService.AdminImpersonateClient(r.Context(), request.GetID(r))
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]string{"token": t}, nil)
}
