package client

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/apikey"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type APIKeyHandler struct {
	apiKeyService *apikey.APIKeyService
}

func NewAPIKeyHandler(apiKeyService *apikey.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{apiKeyService: apiKeyService}
}

func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	keys, err := h.apiKeyService.ListKeys(r.Context(), clientID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve API keys", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, keys, nil)
}

type generateKeyRequest struct {
	Name       string `json:"name"`
	ExpireDays int    `json:"expire_days"`
}

func (h *APIKeyHandler) Generate(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	var req generateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	key, err := h.apiKeyService.GenerateKey(r.Context(), clientID, req.Name, req.ExpireDays)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "GENERATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, key, nil)
}

func (h *APIKeyHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid key ID", nil)
		return
	}

	if err := h.apiKeyService.RevokeKey(r.Context(), id, clientID); err != nil {
		response.Error(w, http.StatusBadRequest, "REVOKE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "revoked",
		"message": "API key revoked successfully",
	}, nil)
}
