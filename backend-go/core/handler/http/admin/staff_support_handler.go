package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func (h *StaffManagementHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	limit, offset := 20, 0
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

	tickets, total, err := h.supportService.ListAllTickets(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve tickets", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, tickets, &response.Meta{Total: total, Limit: limit, Offset: offset})
}

func (h *StaffManagementHandler) ReplyTicket(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	idStr := r.PathValue("id")
	ticketID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ticket ID", nil)
		return
	}

	var req struct {
		Message string `json:"message"`
		Content string `json:"content"` // Support both for compatibility
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	replyText := req.Message
	if replyText == "" {
		replyText = req.Content
	}

	if replyText == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Reply content is required", nil)
		return
	}

	msg, err := h.supportService.StaffReply(r.Context(), ticketID, staffID, replyText)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "REPLY_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, msg, nil)
}
