package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	AgentsKey = "agents"
)

type AgentRepository interface {
	GetCapacity(agentID string) (int, error)
	IncrementCapacity(agentID string) error
	DecrementCapacity(agentID string) error
}

type agentRepository struct {
	client *redis.Client
}

func NewAgentRepository(client *redis.Client) AgentRepository {
	return &agentRepository{
		client: client,
	}
}

// getAgentKey returns Redis key for agent capacity
func (r *agentRepository) getAgentKey(agentID string) string {
	return fmt.Sprintf("%s:%s", AgentsKey, agentID)
}

// GetCapacity gets current number of customers assigned to agent
func (r *agentRepository) GetCapacity(agentID string) (int, error) {
	ctx := context.Background()
	key := r.getAgentKey(agentID)

	// Get current capacity value
	result := r.client.Get(ctx, key)

	// If key doesn't exist, capacity is 0
	if result.Err() == redis.Nil {
		return 0, nil
	}

	if result.Err() != nil {
		return 0, fmt.Errorf("failed to get agent capacity: %w", result.Err())
	}

	// Convert string to int
	capacity, err := strconv.Atoi(result.Val())
	if err != nil {
		return 0, fmt.Errorf("failed to parse capacity value: %w", err)
	}

	return capacity, nil
}

// IncrementCapacity increases agent's customer count by 1
func (r *agentRepository) IncrementCapacity(agentID string) error {
	ctx := context.Background()
	key := r.getAgentKey(agentID)

	// Use INCR to atomically increment the counter
	// If key doesn't exist, it will be set to 1
	err := r.client.Incr(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to increment agent capacity: %w", err)
	}

	return nil
}

// DecrementCapacity decreases agent's customer count by 1
func (r *agentRepository) DecrementCapacity(agentID string) error {
	ctx := context.Background()
	key := r.getAgentKey(agentID)

	// Get current value first to avoid going below 0
	current, err := r.GetCapacity(agentID)
	if err != nil {
		return fmt.Errorf("failed to get current capacity: %w", err)
	}

	// Don't decrement if already at 0
	if current <= 0 {
		return nil
	}

	// Use DECR to atomically decrement the counter
	err = r.client.Decr(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to decrement agent capacity: %w", err)
	}

	return nil
}
