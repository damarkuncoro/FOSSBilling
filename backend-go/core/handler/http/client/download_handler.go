package client

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type DownloadHandler struct {
	downloadService *downloadable.DownloadableService
}

func NewDownloadHandler(downloadService *downloadable.DownloadableService) *DownloadHandler {
	return &DownloadHandler{downloadService: downloadService}
}

func (h *DownloadHandler) GenerateLink(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	fileIDStr := parts[len(parts)-1]
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid file ID", nil)
		return
	}

	link, err := h.downloadService.GenerateDownloadLink(r.Context(), clientID, fileID, 2*time.Hour)
	if err != nil {
		if errors.Is(err, downloadable.ErrProductNotOrdered) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not have an active purchase for this digital good", nil)
			return
		}
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "File not found", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "ACTION_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, link, nil)
}

func (h *DownloadHandler) StreamFile(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	fileIDStr := parts[len(parts)-2]
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid file ID", nil)
		return
	}

	clientIDStr := r.URL.Query().Get("client_id")
	expiresStr := r.URL.Query().Get("expires")
	sig := r.URL.Query().Get("sig")

	clientID, _ := strconv.ParseInt(clientIDStr, 10, 64)
	expiresUnix, _ := strconv.ParseInt(expiresStr, 10, 64)

	file, err := h.downloadService.VerifyAndGetFile(r.Context(), clientID, fileID, expiresUnix, sig)
	if err != nil {
		if errors.Is(err, downloadable.ErrDownloadExpired) {
			response.Error(w, http.StatusGone, "LINK_EXPIRED", "Download link has expired. Please generate a new one.", nil)
			return
		}
		if errors.Is(err, downloadable.ErrInvalidSignature) {
			response.Error(w, http.StatusForbidden, "INVALID_SIGNATURE", "Invalid download signature", nil)
			return
		}
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "File not found", nil)
		return
	}

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", file.Filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("MOCK_BINARY_DATA_FOR_" + file.Filename))
}

func (h *DownloadHandler) ListDownloads(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	downloads, err := h.downloadService.ListClientDownloads(r.Context(), clientID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, downloads, nil)
}

