package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/crnvl96/rate-limiter/infra/cache"
	"github.com/crnvl96/rate-limiter/infra/utilities"
)

type ApiKey struct {
	Key string `json:"key"`
	Id  int64  `json:"id"`
}

type RateLimiterMiddleware struct {
	cache               cache.RateLimiterCache
	requestLimitDefault int64
	requestLimitApiKey  int64
	apiKeys             []ApiKey
}

var GetWd = os.Getwd

func NewRateLimiterMiddleware(cache cache.RateLimiterCache, apiKeyPath string) *RateLimiterMiddleware {
	apiKeys := LoadApiKeys(apiKeyPath)

	limitRequestPerSecondDefault, err := strconv.ParseInt(os.Getenv("LIMIT_REQUEST_PER_SECOND_DEFAULT"), 10, 64)
	if err != nil {
		limitRequestPerSecondDefault = 5
	}

	limitRequestPerSecondApiKey, err := strconv.ParseInt(os.Getenv("LIMITER_REQUEST_PER_SECOND_API_KEY"), 10, 64)
	if err != nil {
		limitRequestPerSecondApiKey = 10
	}

	return &RateLimiterMiddleware{
		cache:               cache,
		requestLimitApiKey:  limitRequestPerSecondApiKey,
		requestLimitDefault: limitRequestPerSecondDefault,
		apiKeys:             apiKeys,
	}
}

func LoadApiKeys(apiKeyPath string) []ApiKey {
	currentDir, err := GetWd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	filePath := filepath.Join(currentDir, apiKeyPath)

	apiKeyJson, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening API key file: %v", err)
	}
	defer apiKeyJson.Close()

	var apiKeys []ApiKey
	if err := json.NewDecoder(apiKeyJson).Decode(&apiKeys); err != nil {
		log.Fatalf("Error decoding API keys: %v", err)
	}

	return apiKeys
}

func (r *RateLimiterMiddleware) findApiKey(apiKey string) (*ApiKey, bool) {
	for _, key := range r.apiKeys {
		if key.Key == apiKey {
			return &key, true
		}
	}
	return nil, false
}

func (r *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip := utilities.GetRealIP(req)
		limiterKeyByIp := "rate-limiter:" + ip
		apiKey := req.Header.Get("API_KEY")

		if apiKey != "" {
			if key, found := r.findApiKey(apiKey); found {
				limiterKeyByApiKey := "rate-limiter:" + strconv.FormatInt(key.Id, 10)
				if r.checkRateLimit(w, limiterKeyByApiKey, r.requestLimitApiKey) {
					return
				}

				go r.cache.ControlExpirationTime(limiterKeyByApiKey)
				next.ServeHTTP(w, req)
				return
			}
		}

		if r.checkRateLimit(w, limiterKeyByIp, r.requestLimitDefault) {
			return
		}

		go r.cache.ControlExpirationTime(limiterKeyByIp)
		next.ServeHTTP(w, req)
	})
}

func (r *RateLimiterMiddleware) checkRateLimit(w http.ResponseWriter, limiterKey string, limit int64) bool {
	limiterKeyAmount, err := r.cache.Get(limiterKey)
	if err != nil && err.Error() != "redis: nil" {
		log.Printf("Error getting rate limit key: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}

	if limiterKeyAmount == "blocked" {
		http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
		return true
	}

	limiterKeyAmountInt, err := strconv.ParseInt(limiterKeyAmount, 10, 64)
	if err != nil {
		limiterKeyAmountInt = 0
	}

	if limiterKeyAmountInt >= limit {
		r.Block(limiterKey)
		http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
		return true
	}

	if err := r.cache.Increment(limiterKey); err != nil {
		log.Printf("Error incrementing rate limit key: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}

	return false
}

func (r *RateLimiterMiddleware) Block(limiterKey string) {
	cacheExpirationStr := os.Getenv("CACHE_EXPIRATION")

	if cacheExpirationStr == "" {
		cacheExpirationStr = "1m" // valor padrão: 1 minuto
	}

	cacheExpiration, err := time.ParseDuration(cacheExpirationStr)
	if err != nil {
		log.Fatalf("Erro ao converter a duração do cache: %v", err)
	}

	r.cache.Set(limiterKey, "blocked", cacheExpiration)
}
