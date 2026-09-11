package guest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	authUc *auth.AuthUsecase
	cfg    *config.Config
	google *oauth2.Config
}

func NewOAuthHandler(au *auth.AuthUsecase, cfg *config.Config) *OAuthHandler {
	g := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  fmt.Sprintf("%s/api/v1/guest/auth/google/callback", cfg.AppURL),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return &OAuthHandler{authUc: au, cfg: cfg, google: g}
}

func (h *OAuthHandler) LoginGoogle(w http.ResponseWriter, r *http.Request) {
	if h.cfg.GoogleClientID == "" {
		response.Error(w, 400, "DISABLED", "Google login is not configured", nil)
		return
	}
	url := h.google.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.google.Exchange(context.Background(), code)
	if err != nil {
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	client := h.google.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Redirect(w, r, "/login?error=fetch_user_failed", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var user struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Redirect(w, r, "/login?error=decode_failed", http.StatusTemporaryRedirect)
		return
	}

	// Login or Register via OAuth
	res, err := h.authUc.OAuthLoginOrRegister(r.Context(), "google", user.ID, user.Email, user.Name)
	if err != nil {
		http.Redirect(w, r, "/login?error="+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	// Redirect back to frontend with token
	target := fmt.Sprintf("/login?token=%s", res.Token)
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
