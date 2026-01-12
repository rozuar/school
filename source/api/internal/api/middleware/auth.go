package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/school-monitoring/backend/internal/auth"
)

// ContextKey tipo para keys de contexto
type ContextKey string

const (
	// UserContextKey key para el usuario en el contexto
	UserContextKey ContextKey = "user"
)

// AuthMiddleware middleware de autenticacion JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		// Verificar formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "Invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validar token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Agregar claims al contexto
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware middleware que verifica el rol del usuario
func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				http.Error(w, `{"error": "User not found in context"}`, http.StatusUnauthorized)
				return
			}

			// Verificar si el rol del usuario esta en los roles permitidos
			roleAllowed := false
			for _, role := range roles {
				if claims.Rol == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				http.Error(w, `{"error": "Insufficient permissions"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// PermissionMiddleware middleware que verifica permisos especificos
func PermissionMiddleware(permisos ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				http.Error(w, `{"error": "User not found in context"}`, http.StatusUnauthorized)
				return
			}

			// Verificar si el usuario tiene al menos uno de los permisos
			if !auth.TieneAlgunPermiso(claims.Rol, permisos...) {
				http.Error(w, `{"error": "Insufficient permissions"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext obtiene los claims del usuario del contexto
func GetUserFromContext(r *http.Request) *auth.Claims {
	claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}
