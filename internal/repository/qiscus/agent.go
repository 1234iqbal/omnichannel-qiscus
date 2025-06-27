package qiscus

import (
	"context"
	"fmt"
	"log"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/pkg/qiscus"
)

type AgentQiscusRepository interface {
	GetAvailableAgents(ctx context.Context) ([]*entity.Agent, error)
	AssignAgentToRoom(ctx context.Context, roomID, agentID string) error
	UnassignAgentFromRoom(ctx context.Context, roomID string) error
}

type agentQiscusRepository struct {
	client *qiscus.Client
}

func NewAgentQiscusRepository(client *qiscus.Client) AgentQiscusRepository {
	return &agentQiscusRepository{
		client: client,
	}
}

func (r *agentQiscusRepository) GetAvailableAgents(ctx context.Context) ([]*entity.Agent, error) {
	log.Printf("Fetching available agents from Qiscus API")

	response, err := r.client.GetAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents from Qiscus: %w", err)
	}

	var agents []*entity.Agent
	var availableCount int

	println("===================")
	println("Processing agents from Qiscus API:")

	for _, qiscusAgent := range response.Data.Agents {
		println(fmt.Sprintf("Agent ID: %d, Name: %s, Email: %s, IsAvailable: %t, ForceOffline: %t, CurrentCustomers: %d",
			qiscusAgent.ID, qiscusAgent.Name, qiscusAgent.Email,
			qiscusAgent.IsAvailable, true))

		if qiscusAgent.IsAvailable {
			agent := &entity.Agent{
				ID:            fmt.Sprintf("%d", qiscusAgent.ID),
				Name:          qiscusAgent.Name,
				Email:         qiscusAgent.Email,
				IsOnline:      true,
				MaxConcurrent: 3,
			}
			agents = append(agents, agent)
			availableCount++

			println(fmt.Sprintf("✅ Added available agent: %s (%s)", agent.Name, agent.Email))
		} else {
			println(fmt.Sprintf("❌ Skipped unavailable agent: %s (Available: %t, ForceOffline: %t)",
				qiscusAgent.Name, qiscusAgent.IsAvailable))
		}
	}

	println("===================")
	log.Printf("Retrieved %d available agents out of %d total agents from Qiscus API",
		availableCount, len(response.Data.Agents))

	return agents, nil
}

func (r *agentQiscusRepository) AssignAgentToRoom(ctx context.Context, roomID, agentID string) error {
	log.Printf("Assigning agent %s to room %s via Qiscus API", agentID, roomID)

	response, err := r.client.AssignAgent(ctx, roomID, agentID)
	if err != nil {
		return fmt.Errorf("failed to assign agent via Qiscus API: %w", err)
	}

	if response.Status != 200 {
		return fmt.Errorf("Qiscus API returned error status %d: %s", response.Status, response.Message)
	}

	log.Printf("Successfully assigned agent %s to room %s", agentID, roomID)
	return nil
}

func (r *agentQiscusRepository) UnassignAgentFromRoom(ctx context.Context, roomID string) error {
	log.Printf("Unassigning agent from room %s via Qiscus API", roomID)

	err := r.client.UnassignAgent(ctx, roomID)
	if err != nil {
		return fmt.Errorf("failed to unassign agent via Qiscus API: %w", err)
	}

	log.Printf("Successfully unassigned agent from room %s", roomID)
	return nil
}
