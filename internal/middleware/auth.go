package middleware

import (
	"context"
	"net/http"

	"github.com/kairos4213/fithub/internal/auth"
	"github.com/kairos4213/fithub/internal/cntx"
	"github.com/kairos4213/fithub/internal/utils"
)

func (mw *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing JWT", err)
			return
		}

		userID, err := auth.ValidateJWT(accessToken, mw.PublicKey)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
			return
		}

		ctx := context.WithValue(r.Context(), cntx.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
