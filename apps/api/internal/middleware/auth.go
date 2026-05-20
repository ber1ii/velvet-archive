package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Define a custom unexported type for context keys to prevent collisions
type contextKey string

const AdminEmailKey contextKey = "admin_email"

// RequireJWT validates the bearer token before allowing the request to hit admin handlers
func RequireJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header missing"}`, http.StatusUnauthorized)
			return
		}

		// Expecting "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "Invalid token format structure"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		// Parse and validate the token claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure token signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Unauthorized: Token is expired or tampered with"}`, http.StatusUnauthorized)
			return
		}

		// Extract custom claims (like email) if needed by handlers later
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if email, ok := claims["email"].(string); ok {
				ctx := context.WithValue(r.Context(), AdminEmailKey, email)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}
