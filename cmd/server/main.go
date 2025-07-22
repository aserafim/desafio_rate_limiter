package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aserafim/desafio_rate_limiter/internal/config"
	"github.com/aserafim/desafio_rate_limiter/internal/limiter"
	"github.com/aserafim/desafio_rate_limiter/internal/limiter/strategy"
)

func main() {
	cfg := config.LoadConfig()

	redisStrategy := strategy.NewRedisLimiter(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Requisição permitida.")
	})

	handler := limiter.RateLimiterMiddleware(cfg, redisStrategy)(mux)

	fmt.Println("Servidor iniciado na porta 8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
