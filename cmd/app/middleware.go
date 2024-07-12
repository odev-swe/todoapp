package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/odev-swe/todoapp/libs"
	"go.uber.org/zap"
)

type UserIdKey string

func (app *application) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		defer func() {
			zap.L().Info("Request Info", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Int("duration - ms", int(time.Since(now).Milliseconds())))
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			libs.Unauthorized(w, "Missing Authorization Header")
			return
		}

		token := strings.Split(tokenHeader, "Bearer ")[1]

		if token == "" {
			libs.Unauthorized(w, "Invalid Token")
			return
		}

		claims, err := libs.ParseToken(token, app.config.JwtSecret)

		if err != nil {
			libs.Unauthorized(w, "Invalid Token")
			return
		}

		id := claims.Data.(map[string]interface{})["id"].(string)

		// wrap in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIdKey("user-id"), id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)

		if allow, retryAfter := app.limiter.Allow(ip); !allow {
			w.Header().Set("Retry-After", fmt.Sprintf("%d s", retryAfter))
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		if err != nil {
			return
		}
		// Rate limit logic
		next.ServeHTTP(w, r)
	})
}
