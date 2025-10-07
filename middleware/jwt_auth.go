package middleware


import (
	"dl/utils"
	"fmt"
	"net/http"
	"strings"
)

// JWTAuth защищает маршруты и добавляет userID в контекст
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Заголовок должен быть формата "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст запроса
		ctx := utils.ContextWithUserID(r.Context(), claims.UserID)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
