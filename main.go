package main

import (
	"log"
	"net/http"

	internal_middleware "github.com/crnvl96/rate-limiter/api/middleware"
	"github.com/crnvl96/rate-limiter/infra/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	cache := cache.NewCache()

	rateLimiter := internal_middleware.NewRateLimiterMiddleware(cache, "api-key.json")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(rateLimiter.Middleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
