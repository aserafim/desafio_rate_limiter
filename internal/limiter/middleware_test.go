package limiter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aserafim/desafio_rate_limiter/internal/config"
)

type mockLimiter struct {
	allow bool
	err   error
}

func (m *mockLimiter) Allow(ctx context.Context, key string, limit, blockDuration int) (bool, error) {
	return m.allow, m.err
}

func TestRateLimiterMiddleware_Allowed(t *testing.T) {
	cfg := &config.Config{
		RateLimitIP:        1,
		BlockDurationIP:    60,
		RateLimitToken:     1,
		BlockDurationToken: 60,
	}

	mock := &mockLimiter{allow: true}
	handler := RateLimiterMiddleware(cfg, mock)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("esperado 200, recebido %d", resp.Code)
	}
}

func TestRateLimiterMiddleware_Blocked(t *testing.T) {
	cfg := &config.Config{
		RateLimitIP:        1,
		BlockDurationIP:    60,
		RateLimitToken:     1,
		BlockDurationToken: 60,
	}

	mock := &mockLimiter{allow: false}
	handler := RateLimiterMiddleware(cfg, mock)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusTooManyRequests {
		t.Errorf("esperado 429, recebido %d", resp.Code)
	}
}
