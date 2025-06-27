package handler

import (
	"encoding/json"
	"log"
	"net/http"

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
	var webhookPayload entity.QiscusWebhookPayload

	if err := json.NewDecoder(r.Body).Decode(&webhookPayload); err != nil {
		log.Printf("Invalid JSON payload: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Incoming webhook for room: %s, user: %s, resolved: %t",
		webhookPayload.RoomID,
		webhookPayload.Email,
		webhookPayload.IsResolved)

	// Convert Qiscus webhook payload to internal RoomLog format
	roomLog := webhookPayload.ToRoomLog()

	// Process based on resolution status
	if webhookPayload.IsResolved {
		// This is a resolved message
		response, err := h.allocationUsecase.ProcessResolvedMessage(r.Context(), roomLog)
		if err != nil {
			log.Printf("Error processing resolved message: %v", err)
			http.Error(w, "Failed to process resolved message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		// This is an incoming message that needs assignment
		response, err := h.allocationUsecase.ProcessIncomingMessage(r.Context(), roomLog)
		if err != nil {
			log.Printf("Error processing incoming message: %v", err)
			http.Error(w, "Failed to process message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}

func (h *WebhookHandler) HandleResolved(w http.ResponseWriter, r *http.Request) {
	var webhookPayload entity.QiscusWebhookPayload

	if err := json.NewDecoder(r.Body).Decode(&webhookPayload); err != nil {
		log.Printf("Invalid JSON payload: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Chat resolved for room: %s", webhookPayload.RoomID)

	// Convert to internal format
	roomLog := webhookPayload.ToRoomLog()

	response, err := h.allocationUsecase.ProcessResolvedMessage(r.Context(), roomLog)
	if err != nil {
		log.Printf("Error processing resolved message: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
