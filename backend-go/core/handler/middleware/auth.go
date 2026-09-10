package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

func RequireAuth(secret string, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := ""; ah := r.Header.Get("Authorization")
			if strings.HasPrefix(strings.ToLower(ah), "bearer ") { t = strings.TrimSpace(ah[7:]) }
			if t == "" { t = r.URL.Query().Get("token") }
			if t == "" { response.Error(w, 401, "ERR", "Missing token", nil); return }
			c, err := auth.ValidateToken(secret, t); if err != nil { response.Error(w, 401, "ERR", "Invalid token: "+err.Error(), nil); return }
			if len(roles) > 0 { ok := false; for _, rl := range roles { if c.Role == rl { ok = true; break } }; if !ok { response.Error(w, 403, "ERR", "Forbidden", nil); return } }
			ctx := context.WithValue(r.Context(), "clientID", c.ClientID)
			ctx = context.WithValue(ctx, "email", c.Email)
			ctx = context.WithValue(ctx, "role", c.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClientID(ctx context.Context) int64 { if v, ok := ctx.Value("clientID").(int64); ok { return v }; return 0 }
func GetRole(ctx context.Context) string { if v, ok := ctx.Value("role").(string); ok { return v }; return "" }
func WithClientID(ctx context.Context, id int64) context.Context { return context.WithValue(ctx, "clientID", id) }
