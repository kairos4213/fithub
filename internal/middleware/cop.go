package middleware

import "net/http"

func (mw *Middleware) Cop(next http.Handler) http.Handler {
	cop := http.NewCrossOriginProtection()
	return cop.Handler(next)
}
