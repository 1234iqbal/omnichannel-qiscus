package usecase

import (
	"encoding/json"
	"fmt"
	"log"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/internal/repository/qiscus"
	"qiscus-agent-allocation/internal/repository/redis"
)

type AllocationUsecase interface {
	IsInQueue(roomID, channel, customerID string) (bool, error)
	AddToQueue(item entity.QueueItem) error
	GetFromQueue() (string, error)

	// Agent operations
	GetOnlineAgents() ([]entity.Agent, error)
	AssignAgent(roomID, agentID string) error
	GetAgentCapacity(agentID string) (int, error)
	IncrementAgentCapacity(agentID string) error
	DecrementAgentCapacity(agentID string) error
}

type allocationUsecase struct {
	agentRepo       redis.AgentRepository
	queueRepo       redis.QueueRepository
	agentQiscusRepo qiscus.AgentQiscusRepository
}

func NewAllocationUsecase(
	agentRepo redis.AgentRepository,
	queueRepo redis.QueueRepository,
	agentQiscusRepo qiscus.AgentQiscusRepository,
) AllocationUsecase {
	return &allocationUsecase{
		agentRepo:       agentRepo,
		queueRepo:       queueRepo,
		agentQiscusRepo: agentQiscusRepo,
	}
}

func (u *allocationUsecase) IsInQueue(roomID, channel, customerID string) (bool, error) {
	exists, err := u.queueRepo.Exists(roomID, channel, customerID)
	if err != nil {
		return false, fmt.Errorf("failed to check if item in queue: %w", err)
	}
	return exists, nil
}

func (u *allocationUsecase) AddToQueue(item entity.QueueItem) error {
	// Convert to JSON string
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal queue item: %w", err)
	}

	// Add to Redis queue (LPUSH for FIFO)
	err = u.queueRepo.Push(string(data))
	if err != nil {
		return fmt.Errorf("failed to push to queue: %w", err)
	}

	log.Printf("Added to queue: %s", string(data))
	return nil
}

// GetFromQueue gets next customer from Redis queue (FIFO)
func (u *allocationUsecase) GetFromQueue() (string, error) {
	// Get from Redis queue (RPOP for FIFO)
	data, err := u.queueRepo.Pop()
	if err != nil {
		return "", fmt.Errorf("failed to pop from queue: %w", err)
	}

	return data, nil
}

// GetOnlineAgents fetches online agents from Qiscus API
func (u *allocationUsecase) GetOnlineAgents() ([]entity.Agent, error) {
	// Get agents from Qiscus API
	qiscusAgents, err := u.agentQiscusRepo.GetOnlineAgents()
	if err != nil {
		return nil, fmt.Errorf("failed to get online agents: %w", err)
	}

	// Convert to our Agent struct
	var agents []entity.Agent
	for _, qAgent := range qiscusAgents {
		agents = append(agents, entity.Agent{
			ID:          fmt.Sprintf("%d", qAgent.ID),
			Name:        qAgent.Name,
			IsAvailable: qAgent.IsAvailable,
		})
	}

	return agents, nil
}

// AssignAgent assigns agent to customer via Qiscus API
func (u *allocationUsecase) AssignAgent(roomID, agentID string) error {
	// Call Qiscus API to assign agent
	err := u.agentQiscusRepo.AssignAgent(roomID, agentID)
	if err != nil {
		return fmt.Errorf("failed to assign agent: %w", err)
	}

	log.Printf("Successfully assigned agent %s to room %s", agentID, roomID)
	return nil
}

// GetAgentCapacity gets current agent capacity from Redis
func (u *allocationUsecase) GetAgentCapacity(agentID string) (int, error) {
	capacity, err := u.agentRepo.GetCapacity(agentID)
	if err != nil {
		return 0, fmt.Errorf("failed to get agent capacity: %w", err)
	}

	return capacity, nil
}

// IncrementAgentCapacity increases agent capacity by 1
func (u *allocationUsecase) IncrementAgentCapacity(agentID string) error {
	err := u.agentRepo.IncrementCapacity(agentID)
	if err != nil {
		return fmt.Errorf("failed to increment agent capacity: %w", err)
	}

	log.Printf("Incremented capacity for agent %s", agentID)
	return nil
}

// DecrementAgentCapacity decreases agent capacity by 1
func (u *allocationUsecase) DecrementAgentCapacity(agentID string) error {
	err := u.agentRepo.DecrementCapacity(agentID)
	if err != nil {
		return fmt.Errorf("failed to decrement agent capacity: %w", err)
	}

	log.Printf("Decremented capacity for agent %s", agentID)
	return nil
}
