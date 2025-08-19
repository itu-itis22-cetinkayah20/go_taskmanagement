package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("gizliAnahtar")

// Context key tipi
type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Token gerekli"}`))
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Geçersiz token"}`))
			return
		}
		// Claims'ten user_id'yi oku ve context'e ekle
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if uid, ok := claims["user_id"].(float64); ok {
				ctx := context.WithValue(r.Context(), UserIDKey, int(uid))
				r = r.WithContext(ctx)
			}
		}
		next(w, r)
	}
}
