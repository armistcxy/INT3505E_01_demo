package middlewares

import (
	"net/http"

	"go.uber.org/ratelimit"
)

func NewRateLimitMiddleware(ratePerSec int) func(next http.Handler) http.Handler {
	rl := ratelimit.New(ratePerSec)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rl.Take()
			next.ServeHTTP(w, r)
		})
	}
}
