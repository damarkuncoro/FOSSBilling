package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/fossbilling/backend-go/internal/handler/middleware"
	"github.com/fossbilling/backend-go/internal/usecase/support"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/fossbilling/backend-go/pkg/response"
)

type SupportHandler struct {
	supportService *support.SupportService
}

func NewSupportHandler(supportService *support.SupportService) *SupportHandler {
	return &SupportHandler{supportService: supportService}
}

type openTicketRequest struct {
	HelpdeskID int64  `json:"helpdesk_id"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
	Priority   string `json:"priority"`
	RelType    *string `json:"rel_type,omitempty"`
	RelID      *int64  `json:"rel_id,omitempty"`
}

func (h *SupportHandler) OpenTicket(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	var req openTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	ticket, err := h.supportService.OpenTicket(r.Context(), support.CreateTicketDTO{
		ClientID:   clientID,
		HelpdeskID: req.HelpdeskID,
		Subject:    req.Subject,
		Message:    req.Message,
		Priority:   "medium",
		RelType:    req.RelType,
		RelID:      req.RelID,
		IPAddress:  r.RemoteAddr,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, ticket, nil)
}

func (h *SupportHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	limit := 20
	offset := 0
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

	tickets, total, err := h.supportService.ListClientTickets(r.Context(), clientID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve tickets", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, tickets, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *SupportHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	ticketIDStr := parts[len(parts)-1]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ticket ID", nil)
		return
	}

	details, err := h.supportService.GetTicket(r.Context(), ticketID, clientID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Ticket not found", nil)
			return
		}
		if errors.Is(err, appErrors.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "You do not have access to this ticket", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get ticket", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, details, nil)
}

type replyRequest struct {
	Message string `json:"message"`
}

func (h *SupportHandler) ReplyTicket(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	ticketIDStr := parts[len(parts)-2]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ticket ID", nil)
		return
	}

	var req replyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	msg, err := h.supportService.ClientReply(r.Context(), ticketID, clientID, req.Message, r.RemoteAddr)
	if err != nil {
		if errors.Is(err, support.ErrTicketClosed) {
			response.Error(w, http.StatusBadRequest, "TICKET_CLOSED", "Cannot reply to a closed ticket", nil)
			return
		}
		if errors.Is(err, appErrors.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "FORBIDDEN", "Forbidden", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "REPLY_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, msg, nil)
}

func (h *SupportHandler) CloseTicket(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetClientID(r.Context())
	if clientID == 0 {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid path", nil)
		return
	}
	ticketIDStr := parts[len(parts)-2]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid ticket ID", nil)
		return
	}

	if err := h.supportService.CloseTicket(r.Context(), ticketID, clientID); err != nil {
		response.Error(w, http.StatusBadRequest, "CLOSE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "closed",
		"message": "Ticket closed successfully",
	}, nil)
}
