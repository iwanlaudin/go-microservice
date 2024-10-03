package api

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type contextKey string

const apiVersionKey contextKey = "api.version"

func ApiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), apiVersionKey, version))
			next.ServeHTTP(w, r)
		})
	}
}

func ErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				log.Printf("Panic: %v\nStack trace: %s", err, debug.Stack())

				NewAppError(nil, "Internal Server Error", http.StatusInternalServerError).SendResponse(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RateLimiter(rps int, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				NewAppError(nil, "Too Many Requests", http.StatusTooManyRequests).SendResponse(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Jika Anda ingin rate limit per IP
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
				NewAppError(nil, "Too Many Requests", http.StatusTooManyRequests).SendResponse(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
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
				NewAppError(nil, "Gateway Timeout", http.StatusGatewayTimeout).SendResponse(w)
				return
			}
		})
	}
}
