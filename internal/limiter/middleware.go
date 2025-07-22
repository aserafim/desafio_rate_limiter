package limiter

import (
	"context"
	"net/http"
	"strings"

	"github.com/aserafim/desafio_rate_limiter/internal/config"
)

func RateLimiterMiddleware(cfg *config.Config, strategy LimiterStrategy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()

			ip := r.RemoteAddr
			token := r.Header.Get("API_KEY")

			var key string
			var limit, block int

			if token != "" {
				key = "token:" + token
				limit = cfg.RateLimitToken
				block = cfg.BlockDurationToken
			} else {
				key = "ip:" + strings.Split(ip, ":")[0]
				limit = cfg.RateLimitIP
				block = cfg.BlockDurationIP
			}

			allowed, err := strategy.Allow(ctx, key, limit, block)
			if err != nil || !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
