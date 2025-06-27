package handler

import (
	"encoding/json"
	"net/http"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type AgentHandler struct {
	agentUsecase usecase.AgentUsecase
}

func NewAgentHandler(agentUsecase usecase.AgentUsecase) *AgentHandler {
	return &AgentHandler{
		agentUsecase: agentUsecase,
	}
}

func (h *AgentHandler) ListAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := h.agentUsecase.GetAllAgents(r.Context())
	if err != nil {
		http.Error(w, "Failed to get agents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (h *AgentHandler) CreateAgent(w http.ResponseWriter, r *http.Request) {
	var agent entity.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.agentUsecase.CreateAgent(r.Context(), &agent)
	if err != nil {
		http.Error(w, "Failed to create agent", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agent)
}

func (h *AgentHandler) UpdateAgentStatus(w http.ResponseWriter, r *http.Request) {
	agentID := chi.URLParam(r, "agentId")

	var statusUpdate entity.AgentStatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	agent, err := h.agentUsecase.UpdateAgentStatus(r.Context(), agentID, &statusUpdate)
	if err != nil {
		http.Error(w, "Agent not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent)
}
