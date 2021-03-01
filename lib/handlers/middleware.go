// Package handlers contains shared router handlers and middleware
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// ContextKey use unique, static context keys to avoid linting warning
// from https://goo.gl/ddu69X
type ContextKey string

func (c ContextKey) String() string {
	return "Context key " + string(c)
}

// ctxKeyTime is the context key for the timestamp at start
var ctxKeyTime = ContextKey("time")

// DefaultMiddlewares sets middleware, MUST BE ADDED BEFORE routes
func DefaultMiddlewares(mux *chi.Mux) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux.Use(addTimeContextMiddleware) // used for request-time and action-time headers
			//r.Use(timeTrackingMiddleware)
			mux.Use(logMiddleware) // own logger
			mux.Use(middleware.RequestID)
			mux.Use(middleware.RealIP)
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.Timeout(180 * time.Second))
		})
	}
}

// logs a request
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		// called at the very end
		defer func() {
			p := r.URL.Path
			if p != "/" && p != "/health" && p != "/ping" {
				log.Info(r.Method, r.URL.Path, "took", time.Since(t1))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// adds the current time to the time context value
// should be added first to the Middleware chain
func addTimeContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKeyTime, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// harmonizes / unifies a given ISO Code
func harmonizeISOCode(code string) string {
	if code == "en" || code == "en-us" || code == "en-gb" {
		code = "en-en"
	}
	if code == "de" || code == "de-at" {
		code = "de-de"
	}
	return code
}
