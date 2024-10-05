package api

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iwanlaudin/go-microservice/pkg/redis"
	"golang.org/x/time/rate"
)

const apiVersionKey contextKey = "api.version"

func ErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				log.Printf("Panic: %v\nStack trace: %s", err, debug.Stack())

				NewAppResponse("Internal Server Error", http.StatusInternalServerError).Err(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func TimeoutMiddleware(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan bool)
			go func() {
				next.ServeHTTP(w, r)
				done <- true
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				NewAppResponse("Gateway Timeout", http.StatusGatewayTimeout).Err(w)
				return
			}
		})
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			NewAppResponse("Missing auth token", http.StatusUnauthorized).Err(w)
			return
		}

		bearerToken := strings.Fields(authHeader)
		if len(bearerToken) < 2 || bearerToken[0] != "Bearer" {
			NewAppResponse("Invalid token format", http.StatusUnauthorized).Err(w)
			return
		}

		user, err := ValidateToken(bearerToken[1])
		if err != nil {
			NewAppResponse(err.Error(), http.StatusUnauthorized).Err(w)
			return
		}

		ctx := ContextWithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ApiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), apiVersionKey, version))
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimiter(rps int, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				NewAppResponse("Too Many Requests", http.StatusTooManyRequests).Err(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

var (
	limiterMap = make(map[string]*rate.Limiter)
	mu         sync.Mutex
)

func getIPLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiterMap[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 5) // 1 request per second, burst of 5
		limiterMap[ip] = limiter
	}

	return limiter
}

func RateLimiterPerIP() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiter := getIPLimiter(ip)
			if !limiter.Allow() {
				NewAppResponse("Too Many Requests", http.StatusTooManyRequests).Err(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RedisRateLimiter(redisClient *redis.RedisClient, limit int, window time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			key := "ratelimit:" + r.RemoteAddr

			count, err := redisClient.Client.Incr(ctx, key).Result()
			if err != nil {
				// Handle error or fallback to allowing the request
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				redisClient.Client.Expire(ctx, key, window)
			}

			if count > int64(limit) {
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
				w.Header().Set("X-RateLimit-Remaining", "0")
				NewAppResponse("Rate limit exceeded", http.StatusTooManyRequests).Err(w)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limit-int(count)))

			next.ServeHTTP(w, r)
		})
	}
}
