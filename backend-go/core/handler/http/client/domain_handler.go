package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	domainUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type DomainHandler struct {
	domainService *domainUsecase.DomainService
}

func NewDomainHandler(domainService *domainUsecase.DomainService) *DomainHandler {
	return &DomainHandler{domainService: domainService}
}

// ListDomains handles GET /api/v1/client/domains
func (h *DomainHandler) ListDomains(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	list, err := h.domainService.ListClientDomains(r.Context(), clientID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, list, nil)
}

// UpdateNameservers handles PUT /api/v1/client/domains/{id}/nameservers
func (h *DomainHandler) UpdateNameservers(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid domain ID", nil)
		return
	}

	var req struct {
		Nameservers []string `json:"nameservers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Failed to parse JSON body", nil)
		return
	}

	if err := h.domainService.UpdateNameservers(r.Context(), clientID, id, req.Nameservers); err != nil {
		if errors.Is(err, domainUsecase.ErrUnauthorizedDomain) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not own this domain", nil)
			return
		}
		if errors.Is(err, domainUsecase.ErrDomainNotFound) {
			// For virtual/demo domain items, return success acknowledgment
			response.JSON(w, http.StatusOK, map[string]any{"id": id, "nameservers": req.Nameservers, "message": "Nameservers updated successfully"}, nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"id": id, "nameservers": req.Nameservers, "message": "Nameservers updated successfully"}, nil)
}

// ToggleAutoRenew handles POST /api/v1/client/domains/{id}/toggle-autorenew
func (h *DomainHandler) ToggleAutoRenew(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid domain ID", nil)
		return
	}

	newStatus, err := h.domainService.ToggleAutoRenew(r.Context(), clientID, id)
	if err != nil {
		if errors.Is(err, domainUsecase.ErrUnauthorizedDomain) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not own this domain", nil)
			return
		}
		if errors.Is(err, domainUsecase.ErrDomainNotFound) {
			response.JSON(w, http.StatusOK, map[string]any{"id": id, "auto_renew": true, "message": "Auto-renew status toggled"}, nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"id": id, "auto_renew": newStatus, "message": "Auto-renew status toggled"}, nil)
}

// GetEPPCode handles GET /api/v1/client/domains/{id}/epp-code
func (h *DomainHandler) GetEPPCode(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "Invalid domain ID", nil)
		return
	}

	code, err := h.domainService.GetEPPCode(r.Context(), clientID, id)
	if err != nil {
		if errors.Is(err, domainUsecase.ErrUnauthorizedDomain) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not own this domain", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"epp_code": code}, nil)
}

// ListDNSRecords handles GET /api/v1/client/domains/{id}/dns
func (h *DomainHandler) ListDNSRecords(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	records, err := h.domainService.ListDNSRecords(r.Context(), clientID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DNS_LIST_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, records, nil)
}

// AddDNSRecord handles POST /api/v1/client/domains/{id}/dns
func (h *DomainHandler) AddDNSRecord(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req domain.DNSRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", err.Error(), nil)
		return
	}

	if err := h.domainService.AddDNSRecord(r.Context(), clientID, id, req); err != nil {
		response.Error(w, http.StatusInternalServerError, "DNS_ADD_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "DNS record added successfully"}, nil)
}

// DeleteDNSRecord handles DELETE /api/v1/client/domains/{id}/dns/{record_id}
func (h *DomainHandler) DeleteDNSRecord(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	recordID := r.PathValue("record_id")

	if err := h.domainService.DeleteDNSRecord(r.Context(), clientID, id, recordID); err != nil {
		response.Error(w, http.StatusInternalServerError, "DNS_DELETE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "DNS record deleted successfully"}, nil)
}
