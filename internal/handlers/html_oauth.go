package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kairos4213/fithub/internal/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func (h *Handler) googleOAuthConfig() *oauth2.Config {
	provider := h.cfg.OAuth["google"]
	return &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthCfg := h.googleOAuthConfig()

	state, err := generateRandomState()
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to generate OAuth state", slog.String("error", err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   600,
	})

	url := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state for CSRF protection
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		h.cfg.Logger.Error("oauth_state cookie missing", slog.String("error", err.Error()))
		http.Error(w, "Invalid OAuth state", http.StatusForbidden)
		return
	}
	if stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "Invalid OAuth state", http.StatusForbidden)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	})

	// Exchange authorization code for tokens
	code := r.URL.Query().Get("code")
	oauthCfg := h.googleOAuthConfig()
	token, err := oauthCfg.Exchange(r.Context(), code)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to exchange OAuth code", slog.String("error", err.Error()))
		return
	}

	// Fetch user info from Google
	gUser, err := fetchGoogleUserInfo(r.Context(), oauthCfg, token)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to fetch Google user info", slog.String("error", err.Error()))
		return
	}

	// Find or create user via account linking
	user, err := h.findOrCreateGoogleUser(r.Context(), gUser)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to find or create Google user", slog.String("error", err.Error()))
		return
	}

	// Issue session tokens (same as password login)
	_, _, err = h.issueSessionTokens(r.Context(), w, user.ID, user.IsAdmin)
	if err != nil {
		HandleInternalServerError(w, r)
		h.cfg.Logger.Error("failed to issue session tokens", slog.String("error", err.Error()))
		return
	}

	if user.IsAdmin {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/workouts", http.StatusSeeOther)
}

func (h *Handler) findOrCreateGoogleUser(ctx context.Context, gUser *googleUserInfo) (database.User, error) {
	provider, err := h.cfg.DB.GetAuthProvider(ctx, database.GetAuthProviderParams{Provider: "google", ProviderUserID: gUser.ID})
	// User exists in auth providers
	if err == nil {
		user, err := h.cfg.DB.GetUserByID(ctx, provider.UserID)
		if err != nil {
			h.cfg.Logger.Error("error fetching user from database", slog.String("error", err.Error()))
			return database.User{}, err
		}

		return user, nil
	}

	user, err := h.cfg.DB.GetUser(ctx, gUser.Email)
	// user does not exist at all
	if err != nil {
		user, err := h.cfg.DB.CreateOAuthUser(ctx, database.CreateOAuthUserParams{
			FirstName:    gUser.GivenName,
			LastName:     gUser.FamilyName,
			Email:        gUser.Email,
			ProfileImage: sql.NullString{String: gUser.Picture, Valid: true},
		})
		if err != nil {
			h.cfg.Logger.Error("error creating oauth user", slog.String("error", err.Error()))
			return database.User{}, err
		}

		_, err = h.cfg.DB.CreateAuthProvider(ctx, database.CreateAuthProviderParams{
			UserID:         user.ID,
			Provider:       "google",
			ProviderUserID: gUser.ID,
		})
		if err != nil {
			h.cfg.Logger.Error("error linking user to oauth provider", "error", err.Error())
			return database.User{}, err
		}

		return user, nil
	}

	// user exists
	_, err = h.cfg.DB.CreateAuthProvider(ctx, database.CreateAuthProviderParams{
		UserID:         user.ID,
		Provider:       "google",
		ProviderUserID: gUser.ID,
	})
	if err != nil {
		h.cfg.Logger.Error("error linking user to oauth provider", "error", err.Error())
		return database.User{}, err
	}

	return user, nil
}

func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*googleUserInfo, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
