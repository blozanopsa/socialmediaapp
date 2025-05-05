package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware checks for a Bearer token (Microsoft OAuth) in the Authorization header
// or a valid sessionID cookie
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		sessionCookie, err := r.Cookie("sessionID")
		if (authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ")) && (err != nil || sessionCookie.Value == "") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// TODO: Validate the Microsoft OAuth token here if using Bearer
		// TODO: Validate the sessionID cookie here if using session
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers to allow cross-origin requests
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure CORS headers are applied to all responses
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
