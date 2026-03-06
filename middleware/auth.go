package middleware

import (
    "context"
    "net/http"
    "os"
    "strings"
	"encoding/json"

    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        // ambil header Authorization
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            writeJSON(w, 401, map[string]any{"status": "error", "message": "Token tidak ada"})
            return
        }

        // cek format "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            writeJSON(w, 401, map[string]any{"status": "error", "message": "Format token tidak valid"})
            return
        }

        tokenString := parts[1]

        // parse & validasi token
        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            writeJSON(w, 401, map[string]any{"status": "error", "message": "Token tidak valid atau expired"})
            return
        }

        // ambil userID dari token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            writeJSON(w, 401, map[string]any{"status": "error", "message": "Token tidak valid"})
            return
        }

        userID := int(claims["userID"].(float64))

        // simpan userID ke context → bisa dipakai di handler
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func writeJSON(w http.ResponseWriter, status int, data map[string]any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}