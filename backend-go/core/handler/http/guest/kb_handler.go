package guest

import (
	"errors"
	"net/http"
	"strconv"

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

func (h *KBHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.ListCategories(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, cats, nil)
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

func (h *KBHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Slug is required", nil)
		return
	}

	art, err := h.svc.GetArticle(r.Context(), slug)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Article not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}

	response.JSON(w, http.StatusOK, art, nil)
}
