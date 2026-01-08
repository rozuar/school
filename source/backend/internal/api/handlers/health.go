package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/school-monitoring/backend/internal/api/middleware"
)

// Health responde para healthchecks (Railway / LB)
func Health(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":         true,
		"service":    "backend",
		"request_id": middleware.GetRequestID(r),
		"time":       time.Now().UTC().Format(time.RFC3339),
	})
}


