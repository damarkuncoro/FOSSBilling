package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/news"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/staff"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type NewsHandler struct {
	newsService  *news.NewsService
	staffService *staff.StaffService
}

func NewNewsHandler(newsService *news.NewsService, staffService *staff.StaffService) *NewsHandler {
	return &NewsHandler{
		newsService:  newsService,
		staffService: staffService,
	}
}

func (h *NewsHandler) List(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "news", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: news", nil)
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

	posts, total, err := h.newsService.ListAll(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve news", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *NewsHandler) Create(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "news", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: news", nil)
		return
	}

	var dto news.CreateNewsDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}
	dto.AdminID = staffID

	post, err := h.newsService.Create(r.Context(), dto)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "CREATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, post, nil)
}

func (h *NewsHandler) Update(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "news", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: news", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid news ID", nil)
		return
	}

	var dto news.UpdateNewsDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", nil)
		return
	}

	post, err := h.newsService.Update(r.Context(), id, dto)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "News article not found", nil)
			return
		}
		response.Error(w, http.StatusBadRequest, "UPDATE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, post, nil)
}

func (h *NewsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "news", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: news", nil)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid news ID", nil)
		return
	}

	if err := h.newsService.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, "DELETE_FAILED", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "deleted",
		"message": "News article deleted successfully",
	}, nil)
}
