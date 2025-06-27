package qiscus

import (
	"fmt"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/pkg/qiscus"
)

type AgentQiscusRepository interface {
	GetOnlineAgents() ([]entity.QiscusAgent, error)
	AssignAgent(roomID, agentID string) error
}

type agentQiscusRepository struct {
	client *qiscus.Client
}

func NewAgentQiscusRepository(client *qiscus.Client) AgentQiscusRepository {
	return &agentQiscusRepository{
		client: client,
	}
}

// GetOnlineAgents fetches online agents from Qiscus API
func (r *agentQiscusRepository) GetOnlineAgents() ([]entity.QiscusAgent, error) {
	// Call Qiscus API to get agents
	agents, err := r.client.GetAgents()
	if err != nil {
		return nil, fmt.Errorf("failed to get agents from Qiscus: %w", err)
	}

	// Filter only online agents
	var onlineAgents []entity.QiscusAgent
	for _, agent := range agents {
		if agent.IsAvailable {
			onlineAgents = append(onlineAgents, entity.QiscusAgent{
				ID:          agent.ID,
				Name:        agent.Name,
				IsAvailable: agent.IsAvailable,
			})
		}
	}

	return onlineAgents, nil
}

// AssignAgent assigns an agent to a room via Qiscus API
func (r *agentQiscusRepository) AssignAgent(roomID, agentID string) error {
	// Call Qiscus API to assign agent
	err := r.client.AssignAgent(roomID, agentID)
	if err != nil {
		return fmt.Errorf("failed to assign agent via Qiscus API: %w", err)
	}

	return nil
}
