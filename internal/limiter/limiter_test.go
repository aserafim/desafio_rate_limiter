package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"

	"github.com/aserafim/desafio_rate_limiter/internal/config"
	"github.com/aserafim/desafio_rate_limiter/internal/limiter/strategy"
)

func init() {
	_ = godotenv.Load("../.env")
}

func TestRateLimiterByIP(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cfg := &config.Config{
		RateLimitIP:        2,
		BlockDurationIP:    3, // seconds
		RateLimitToken:     0, // Not used in this test
		BlockDurationToken: 0, // Not used in this test
	}

	redisStrategy := strategy.NewRedisLimiter(rdb.Options().Addr, rdb.Options().Password, rdb.Options().DB)

	mw := RateLimiterMiddleware(cfg, redisStrategy)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	// Primeira requisição - deve passar
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}

	// Segunda requisição - deve passar
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}

	// Terceira requisição - deve ser bloqueada
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("esperado 429, obtido %d", rec.Code)
	}
}

func TestRateLimiterByToken(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cfg := &config.Config{
		RateLimitIP:        0, // Not used in this test
		BlockDurationIP:    0, // Not used in this test
		RateLimitToken:     1,
		BlockDurationToken: 2, // seconds
	}

	redisStrategy := strategy.NewRedisLimiter(rdb.Options().Addr, rdb.Options().Password, rdb.Options().DB)

	mw := RateLimiterMiddleware(cfg, redisStrategy)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "abc123")
	rec := httptest.NewRecorder()

	// Primeira requisição - deve passar
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtido %d", rec.Code)
	}

	// Segunda requisição - deve ser bloqueada
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("esperado 429, obtido %d", rec.Code)
	}

	// Aguarda expiração
	t.Log("Aguardando expiração...")
	time.Sleep(3 * time.Second)

	// Terceira requisição - deve passar após expiração
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperado 200 após expiração, obtido %d", rec.Code)
	}
}

func TestTokenOverridesIP(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cfg := &config.Config{
		RateLimitIP:        1,
		BlockDurationIP:    2, // seconds
		RateLimitToken:     3,
		BlockDurationToken: 2, // seconds
	}

	redisStrategy := strategy.NewRedisLimiter(rdb.Options().Addr, rdb.Options().Password, rdb.Options().DB)

	mw := RateLimiterMiddleware(cfg, redisStrategy)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "vip-token")

	// Envia 3 requisições com token especial
	for i := 0; i < 3; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("esperado 200 na %dª, obtido %d", i+1, rec.Code)
		}
	}

	// 4ª deve falhar
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("esperado 429 após 4 requisições, obtido %d", rec.Code)
	}
}
