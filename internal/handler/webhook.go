package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/internal/usecase"
)

type WebhookHandler struct {
	allocationUsecase usecase.AllocationUsecase
}

func NewWebhookHandler(allocationUsecase usecase.AllocationUsecase) *WebhookHandler {
	return &WebhookHandler{
		allocationUsecase: allocationUsecase,
	}
}

func (h *WebhookHandler) HandleIncoming(w http.ResponseWriter, r *http.Request) {
	var webhook entity.QiscusWebhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		log.Printf("Failed to decode webhook: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. Validate webhook payload
	if webhook.RoomID == "" || webhook.Email == "" {
		log.Println("Invalid webhook payload: missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Skip if already resolved
	if webhook.IsResolved {
		log.Printf("Chat already resolved, skipping: Room %s", webhook.RoomID)
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. Check if already in queue (prevent duplicate)
	exists, err := h.allocationUsecase.IsInQueue(webhook.RoomID, webhook.Source, webhook.Email)
	if err != nil {
		log.Printf("Failed to check queue: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Printf("Room already in queue, skipping: Room %s", webhook.RoomID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "already_queued",
			"message": "Customer already in queue",
		})
		return
	}

	// 4. Add to Redis Queue (FIFO with timestamp)
	queueItem := entity.QueueItem{
		CustomerID: webhook.Email,
		RoomID:     webhook.RoomID,
		Channel:    webhook.Source,
		Timestamp:  time.Now(),
	}

	err = h.allocationUsecase.AddToQueue(queueItem)
	if err != nil {
		log.Printf("Failed to add to queue: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "queued",
		"message": "Customer added to queue successfully",
	})
}

func (h *WebhookHandler) HandleResolved(w http.ResponseWriter, r *http.Request) {
	var webhook entity.QiscusWebhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		log.Printf("Failed to decode webhook: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 2. Validate webhook payload
	if webhook.RoomID == "" {
		log.Println("Invalid webhook payload: missing room_id")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Skip if not resolved
	if !webhook.IsResolved {
		log.Printf("Chat not resolved yet, skipping: Room %s", webhook.RoomID)
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. Update Agent Capacity (Redis -1) if agent exists
	if webhook.CandidateAgent.ID > 0 {
		agentID := fmt.Sprintf("%d", webhook.CandidateAgent.ID)
		err := h.allocationUsecase.DecrementAgentCapacity(agentID)
		if err != nil {
			log.Printf("Failed to decrement agent capacity: %v", err)
			// Continue processing even if this fails
		}
		log.Printf("Chat resolved: Room %s, Agent %s", webhook.RoomID, agentID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "resolved",
		"message": "Chat resolved successfully",
	})
}
