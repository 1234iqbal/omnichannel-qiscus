package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"qiscus-agent-allocation/internal/domain/entity"

	"github.com/go-redis/redis/v8"
)

const (
	AgentsKey        = "agents"
	AgentChatsPrefix = "agent_chats:"
)

type AgentRepository interface {
	Create(ctx context.Context, agent *entity.Agent) error
	GetByID(ctx context.Context, agentID string) (*entity.Agent, error)
	GetAll(ctx context.Context) ([]*entity.Agent, error)
	Update(ctx context.Context, agent *entity.Agent) error
	Delete(ctx context.Context, agentID string) error

	FindAgentByRoomID(ctx context.Context, roomID string) (string, error)
}

type agentRepository struct {
	client *redis.Client
}

func NewAgentRepository(client *redis.Client) AgentRepository {
	return &agentRepository{
		client: client,
	}
}

func (r *agentRepository) Create(ctx context.Context, agent *entity.Agent) error {
	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	return r.client.HSet(ctx, AgentsKey, agent.ID, agentJSON).Err()
}

func (r *agentRepository) GetByID(ctx context.Context, agentID string) (*entity.Agent, error) {
	agentJSON, err := r.client.HGet(ctx, AgentsKey, agentID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("agent not found")
		}
		return nil, err
	}

	var agent entity.Agent
	err = json.Unmarshal([]byte(agentJSON), &agent)
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (r *agentRepository) GetAll(ctx context.Context) ([]*entity.Agent, error) {
	agentsMap, err := r.client.HGetAll(ctx, AgentsKey).Result()
	if err != nil {
		return nil, err
	}

	var agents []*entity.Agent
	for _, agentJSON := range agentsMap {
		var agent entity.Agent
		if err := json.Unmarshal([]byte(agentJSON), &agent); err != nil {
			continue
		}
		agents = append(agents, &agent)
	}

	return agents, nil
}

func (r *agentRepository) Update(ctx context.Context, agent *entity.Agent) error {
	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	return r.client.HSet(ctx, AgentsKey, agent.ID, agentJSON).Err()
}

func (r *agentRepository) Delete(ctx context.Context, agentID string) error {
	// Remove from agents hash
	err := r.client.HDel(ctx, AgentsKey, agentID).Err()
	if err != nil {
		return err
	}

	// Remove agent's chat set
	return r.client.Del(ctx, AgentChatsPrefix+agentID).Err()
}

func (r *agentRepository) FindAgentByRoomID(ctx context.Context, roomID string) (string, error) {
	agents, err := r.GetAll(ctx)
	if err != nil {
		return "", err
	}

	for _, agent := range agents {

		isMember, err := r.client.SIsMember(ctx, AgentChatsPrefix+agent.ID, roomID).Result()
		if err != nil {
			println(fmt.Sprintf("Error checking agent %s: %v", agent.ID, err))
			continue // Skip this agent if there's an error
		}

		println(fmt.Sprintf("Agent %s has room %s: %t", agent.ID, roomID, isMember))

		if isMember {
			println(fmt.Sprintf("Found! Agent %s is handling room %s", agent.ID, roomID))
			return agent.ID, nil // Found the agent handling this room
		}
	}

	println(fmt.Sprintf("No agent found handling room %s", roomID))
	return "", fmt.Errorf("no agent found for room %s", roomID)
}
