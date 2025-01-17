package handlers

import (
	"fmt"
	"forum/internal/cookies"
	"net/http"
	"sync"
	"time"
)

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				h.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) CheckGuest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie := cookies.GetSessionCookie("session_id", r)
		if sessionCookie != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RequireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie := cookies.GetSessionCookie("session_id", r)
		if sessionCookie == nil || !h.service.IsSessionValid(sessionCookie.Value) {
			if sessionCookie != nil {
				cookies.ExpireSessionCookie("session_id", w)
			}
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RequireModerator(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie := cookies.GetSessionCookie("session_id", r)
		if sessionCookie == nil || !h.service.IsSessionValid(sessionCookie.Value) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		user, err := h.service.GetUser(r)
		if err != nil || user.Role != "Moderator" {
			http.Error(w, "Forbidden: Moderator access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie := cookies.GetSessionCookie("session_id", r)
		if sessionCookie == nil || !h.service.IsSessionValid(sessionCookie.Value) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		user, err := h.service.GetUser(r)
		if err != nil || user.Role != "Admin" {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type client struct {
	limiter  *rateLimiter
	lastSeen time.Time
}

type rateLimiter struct {
	mu        sync.Mutex
	tokens    int
	lastCheck time.Time
}

func newRateLimiter(maxTokens int, refillRate time.Duration) *rateLimiter {
	return &rateLimiter{
		tokens:    maxTokens,
		lastCheck: time.Now(),
	}
}

func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck)
	refillTokens := int(elapsed / time.Second)
	rl.tokens += refillTokens
	if rl.tokens > 5 {
		rl.tokens = 5
	}
	rl.lastCheck = now

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func (h *Handler) RateLimiter(next http.Handler) http.Handler {
	clients := make(map[string]*client)
	mu := sync.Mutex{}

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 2*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()
		c, exists := clients[ip]
		if !exists {
			c = &client{
				limiter:  newRateLimiter(5, time.Second),
				lastSeen: time.Now(),
			}
			clients[ip] = c
		}
		c.lastSeen = time.Now()
		mu.Unlock()

		if !c.limiter.allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
