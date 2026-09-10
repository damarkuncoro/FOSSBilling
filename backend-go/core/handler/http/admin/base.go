package admin

import (
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func check(w http.ResponseWriter, r *http.Request, s *staff.StaffService, mod, act string) bool {
	ok, _ := s.HasPermission(r.Context(), middleware.GetClientID(r.Context()), mod, act)
	if !ok {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Missing permission: "+mod+":"+act, nil)
	}
	return ok
}
