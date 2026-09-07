package client

import (
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/notification"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct {
	notifService *notification.NotificationService
}

func NewNotificationHandler(notifService *notification.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notifService: notifService,
	}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	notifications, total, err := h.notifService.ListMyNotifications(r.Context(), clientID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	meta := &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	response.JSON(w, http.StatusOK, notifications, meta)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id == 0 {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "invalid notification id", nil)
		return
	}

	if err := h.notifService.MarkAsRead(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())

	if err := h.notifService.MarkAllAsRead(r.Context(), clientID); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}
