package main

import (
	"fmt"
	"log"
	"net/http"

	"qiscus-agent-allocation/internal/config"
	"qiscus-agent-allocation/internal/handler"
	qiscusRepo "qiscus-agent-allocation/internal/repository/qiscus"
	redisRepo "qiscus-agent-allocation/internal/repository/redis"
	"qiscus-agent-allocation/internal/usecase"
	"qiscus-agent-allocation/pkg/qiscus"
	redisClient "qiscus-agent-allocation/pkg/redis"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Redis client
	client, err := redisClient.NewClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	defer client.Close()

	qiscusClient := qiscus.NewClient(qiscus.Config{
		BaseURL:   cfg.QiscusConfig.BaseURL,
		AppID:     cfg.QiscusConfig.AppID,
		SecretKey: cfg.QiscusConfig.SecretKey,
		Timeout:   cfg.QiscusConfig.Timeout,
	})

	// Initialize repositories
	agentRepo := redisRepo.NewAgentRepository(client)
	queueRepo := redisRepo.NewQueueRepository(client)

	agentQiscusRepo := qiscusRepo.NewAgentQiscusRepository(qiscusClient)

	// Initialize use cases
	allocationUsecase := usecase.NewAllocationUsecase(agentRepo, queueRepo, agentQiscusRepo)

	// Initialize handlers
	webhookHandler := handler.NewWebhookHandler(allocationUsecase)

	// Setup routes
	r := chi.NewRouter()

	// Routes
	r.Route("/webhook", func(r chi.Router) {
		r.Post("/incoming", webhookHandler.HandleIncoming)
		r.Post("/resolved", webhookHandler.HandleResolved)
	})

	// start worker

	fmt.Printf("Server starting on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
