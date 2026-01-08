package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDContextKey ContextKey = "request_id"

// RequestIDMiddleware asegura un request id para correlación (logs/headers).
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		w.Header().Set("X-Request-Id", reqID)
		ctx := context.WithValue(r.Context(), RequestIDContextKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


