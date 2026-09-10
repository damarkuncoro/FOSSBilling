package admin

import (
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/knowledgebase"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type KBHandler struct {
	staffService *staff.StaffService; svc *knowledgebase.Service
}

func NewKBHandler(s *staff.StaffService, svc *knowledgebase.Service) *KBHandler { return &KBHandler{s, svc} }

func (h *KBHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "knowledgebase", "read") { return }
	cid, _ := strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64)
	l, o := request.GetLimitOffset(r)
	as, tot, err := h.svc.ListArticles(r.Context(), cid, l, o)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, as, &response.Meta{Total: tot, Limit: l, Offset: o})
}

func (h *KBHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "knowledgebase", "write") { return }
	var a domain.KBArticle; if request.Decode(r, &a) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.svc.CreateArticle(r.Context(), &a); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, a, nil)
}

func (h *KBHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "knowledgebase", "write") { return }
	var a domain.KBArticle; if request.Decode(r, &a) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	a.ID = request.GetID(r)
	if err := h.svc.UpdateArticle(r.Context(), &a); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, a, nil)
}

func (h *KBHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	if !check(w, r, h.staffService, "knowledgebase", "delete") { return }
	if err := h.svc.DeleteArticle(r.Context(), request.GetID(r)); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, map[string]bool{"success": true}, nil)
}

func (h *KBHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cs, err := h.svc.ListCategories(r.Context())
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, cs, nil)
}

func (h *KBHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var c domain.KBCategory; if request.Decode(r, &c) != nil { response.Error(w, 400, "BAD", "Invalid", nil); return }
	if err := h.svc.CreateCategory(r.Context(), &c); err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 201, c, nil)
}
