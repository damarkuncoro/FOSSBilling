package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	pkgAuth "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

type ClientManagementHandler struct {
	staffService *staff.StaffService
	authService  *auth.AuthUsecase
	clientRepo   domain.ClientRepository
}

func NewClientManagementHandler(staffService *staff.StaffService, authService *auth.AuthUsecase, clientRepo domain.ClientRepository) *ClientManagementHandler {
	return &ClientManagementHandler{staffService: staffService, authService: authService, clientRepo: clientRepo}
}

func (h *ClientManagementHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "clients", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: clients", nil)
		return
	}

	limit, offset := 50, 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	clients, total, err := h.clientRepo.List(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve clients", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, clients, &response.Meta{Total: total, Limit: limit, Offset: offset})
}

func (h *ClientManagementHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "clients", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: clients", nil)
		return
	}

	idStr := r.PathValue("id")
	clientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid client ID", nil)
		return
	}

	client, err := h.clientRepo.GetByID(r.Context(), clientID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Client not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get client", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, client, nil)
}

type ClientPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Company   string `json:"company"`
	Country   string `json:"country"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
}

func (h *ClientManagementHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "clients", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: clients", nil)
		return
	}

	var req ClientPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", nil)
		return
	}
	if strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.Email) == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_FAILED", "First name and email are required", nil)
		return
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}
	if req.Status == "" {
		req.Status = "active"
	}
	pwd := req.Password
	if pwd == "" {
		pwd = "Password123!"
	}
	passHash, _ := pkgAuth.HashPassword(pwd)

	client := &domain.Client{
		Email: req.Email, PasswordHash: passHash,
		FirstName: security.SanitizeAlphaNumeric(req.FirstName),
		LastName:  security.SanitizeAlphaNumeric(req.LastName),
		Company:   security.SanitizeHTML(req.Company),
		Country:   security.SanitizeAlphaNumeric(req.Country),
		Currency:  security.SanitizeAlphaNumeric(req.Currency),
		Status:    domain.ClientStatus(req.Status),
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}

	if err := h.clientRepo.Create(r.Context(), client); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save client", err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, client, nil)
}

func (h *ClientManagementHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "clients", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: clients", nil)
		return
	}

	idStr := r.PathValue("id")
	clientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid client ID", nil)
		return
	}

	client, err := h.clientRepo.GetByID(r.Context(), clientID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Client not found", nil)
		return
	}

	var req ClientPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body", nil)
		return
	}

	if req.FirstName != "" {
		client.FirstName = security.SanitizeAlphaNumeric(req.FirstName)
	}
	if req.LastName != "" {
		client.LastName = security.SanitizeAlphaNumeric(req.LastName)
	}
	if req.Email != "" {
		client.Email = req.Email
	}
	if req.Company != "" {
		client.Company = security.SanitizeHTML(req.Company)
	}
	if req.Country != "" {
		client.Country = security.SanitizeAlphaNumeric(req.Country)
	}
	if req.Currency != "" {
		client.Currency = security.SanitizeAlphaNumeric(req.Currency)
	}
	if req.Status != "" {
		client.Status = domain.ClientStatus(req.Status)
	}
	if req.Password != "" {
		client.PasswordHash, _ = pkgAuth.HashPassword(req.Password)
	}

	if err := h.clientRepo.Update(r.Context(), client); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update client", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, client, nil)
}

func (h *ClientManagementHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "clients", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: clients", nil)
		return
	}

	idStr := r.PathValue("id")
	clientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid client ID", nil)
		return
	}

	if err := h.clientRepo.Delete(r.Context(), clientID); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Client not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete client", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

func (h *ClientManagementHandler) ImpersonateClient(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "clients", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: clients", nil)
		return
	}

	idStr := r.PathValue("id")
	clientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid client ID", nil)
		return
	}

	token, err := h.authService.AdminImpersonateClient(r.Context(), clientID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Impersonation failed", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"token": token}, nil)
}
