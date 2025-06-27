package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type HealthHandler struct {
	redisClient *redis.Client
}

func NewHealthHandler(redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		redisClient: redisClient,
	}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check Redis connection
	_, err := h.redisClient.Ping(ctx).Result()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  "Redis connection failed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"services": map[string]string{
			"redis": "connected",
		},
	})
}
