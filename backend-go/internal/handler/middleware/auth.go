package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type contextKey string

const (
	ClientIDKey contextKey = "clientID"
	EmailKey    contextKey = "email"
	RoleKey     contextKey = "role"
)

func RequireAuth(jwtSecret string, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header is required", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header must be Bearer <token>", nil)
				return
			}

			claims, err := auth.ValidateToken(jwtSecret, parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "INVALID_TOKEN", "Token is invalid or expired", nil)
				return
			}

			if len(allowedRoles) > 0 {
				roleAllowed := false
				for _, r := range allowedRoles {
					if claims.Role == r {
						roleAllowed = true
						break
					}
				}
				if !roleAllowed {
					response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient role permissions", nil)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ClientIDKey, claims.ClientID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClientID(ctx context.Context) int64 {
	if val, ok := ctx.Value(ClientIDKey).(int64); ok {
		return val
	}
	return 0
}
