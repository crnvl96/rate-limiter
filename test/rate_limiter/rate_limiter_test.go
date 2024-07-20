package rate_limiter_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/crnvl96/rate-limiter/api/middleware"
	"github.com/crnvl96/rate-limiter/infra/cache"
)

func TestRateLimiterMiddlewareByIP(t *testing.T) {
	cache := &cache.MockCache{
		Data: make(map[string]string),
	}

	os.Setenv("LIMIT_REQUEST_PER_SECOND_DEFAULT", "5")
	os.Setenv("LIMITER_REQUEST_PER_SECOND_API_KEY", "10")

	middleware := middleware.NewRateLimiterMiddleware(cache, "./test_api_keys.json")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "")

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}
}

func TestRateLimiterMiddlewareByAPIKey(t *testing.T) {
	cache := &cache.MockCache{
		Data: make(map[string]string),
	}

	os.Setenv("LIMIT_REQUEST_PER_SECOND_DEFAULT", "10")
	os.Setenv("LIMITER_REQUEST_PER_SECOND_API_KEY", "10")

	middleware := middleware.NewRateLimiterMiddleware(cache, "./test_api_keys.json")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "key1")

	for i := 0; i < 10; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}
}

func TestRateLimiterMiddlewareByAPIKeyNotFound(t *testing.T) {
	cache := &cache.MockCache{Data: make(map[string]string)}
	os.Setenv("LIMIT_REQUEST_PER_SECOND_DEFAULT", "5")

	middleware := middleware.NewRateLimiterMiddleware(cache, "./test_api_keys.json")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "notfound")

	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}
}

func TestRateLimiterIfBlockTheUserWhenTheLimitIsReached(t *testing.T) {
	cache := &cache.MockCache{
		Data: make(map[string]string),
	}
	os.Setenv("LIMIT_REQUEST_PER_SECOND_DEFAULT", "10")
	os.Setenv("LIMITER_REQUEST_PER_SECOND_API_KEY", "10")

	middleware := middleware.NewRateLimiterMiddleware(cache, "./test_api_keys.json")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "key1")

	for i := 0; i < 10; i++ {
		rr := httptest.NewRecorder()
		middleware.Middleware(handler).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	middleware.Middleware(handler).ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
	}

	limiter, _ := cache.Get("rate-limiter:" + strconv.FormatInt(1, 10))

	if limiter != "blocked" {
		t.Errorf("Expected blocked, got %v", limiter)
	}
}
