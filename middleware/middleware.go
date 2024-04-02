package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, duration)
	}
}

// Auth Middleware that checks cookie session
func Auth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if cookie.Value != "123" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func JSONMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}

// Middleware is a function that receives a http.Handler and returns a http.Handler
type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(m ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next.ServeHTTP
	}
}
