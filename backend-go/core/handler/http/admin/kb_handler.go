package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/knowledgebase"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type KBHandler struct {
	svc *knowledgebase.Service
}

func NewKBHandler(svc *knowledgebase.Service) *KBHandler {
	return &KBHandler{svc: svc}
}

func (h *KBHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	catID, _ := strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	arts, total, err := h.svc.ListArticles(r.Context(), catID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, arts, &response.Meta{Total: total, Limit: limit, Offset: offset})
}

func (h *KBHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var art domain.KBArticle
	if err := json.NewDecoder(r.Body).Decode(&art); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}

	if err := h.svc.CreateArticle(r.Context(), &art); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, art, nil)
}

func (h *KBHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var art domain.KBArticle
	if err := json.NewDecoder(r.Body).Decode(&art); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}
	art.ID = id

	if err := h.svc.UpdateArticle(r.Context(), &art); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Article not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, art, nil)
}

func (h *KBHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteArticle(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"success": true}, nil)
}

func (h *KBHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.ListCategories(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, cats, nil)
}

func (h *KBHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat domain.KBCategory
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}

	if err := h.svc.CreateCategory(r.Context(), &cat); err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusCreated, cat, nil)
}
