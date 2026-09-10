package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/downloadable"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/request"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type DownloadHandler struct{ svc *downloadable.DownloadableService }

func NewDownloadHandler(s *downloadable.DownloadableService) *DownloadHandler { return &DownloadHandler{s} }

func (h *DownloadHandler) GenerateLink(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	l, err := h.svc.GenerateDownloadLink(r.Context(), cid, request.GetID(r), 2*time.Hour)
	if err != nil { response.Error(w, 403, "FORBIDDEN", err.Error(), nil); return }
	response.JSON(w, 200, l, nil)
}

func (h *DownloadHandler) StreamFile(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseInt(r.URL.Query().Get("client_id"), 10, 64)
	exp, _ := strconv.ParseInt(r.URL.Query().Get("expires"), 10, 64)
	f, err := h.svc.VerifyAndGetFile(r.Context(), cid, request.GetID(r), exp, r.URL.Query().Get("sig"))
	if err != nil { response.Error(w, 403, "FORBIDDEN", err.Error(), nil); return }
	w.Header().Set("Content-Type", f.ContentType); w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", f.Filename))
	w.Write([]byte("MOCK_BINARY_DATA"))
}

func (h *DownloadHandler) ListDownloads(w http.ResponseWriter, r *http.Request) {
	cid := middleware.GetClientID(r.Context()); if cid == 0 { response.Error(w, 401, "UNAUTHORIZED", "Login required", nil); return }
	ds, err := h.svc.ListClientDownloads(r.Context(), cid)
	if err != nil { response.Error(w, 500, "ERR", err.Error(), nil); return }
	response.JSON(w, 200, ds, nil)
}
