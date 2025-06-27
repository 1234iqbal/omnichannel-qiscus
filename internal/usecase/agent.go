package usecase

import (
	"context"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/internal/repository/qiscus"
	"qiscus-agent-allocation/internal/repository/redis"
)

const DefaultMaxChats = 2

type AgentUsecase interface {
	CreateAgent(ctx context.Context, agent *entity.Agent) error
	GetAgent(ctx context.Context, agentID string) (*entity.Agent, error)
	GetAllAgents(ctx context.Context) ([]*entity.Agent, error)
	UpdateAgentStatus(ctx context.Context, agentID string, update *entity.AgentStatusUpdate) (*entity.Agent, error)
}

type agentUsecase struct {
	agentRepo       redis.AgentRepository
	agentQiscusRepo qiscus.AgentQiscusRepository
}

func NewAgentUsecase(agentRepo redis.AgentRepository, agentQiscusRepo qiscus.AgentQiscusRepository) AgentUsecase {
	return &agentUsecase{
		agentRepo:       agentRepo,
		agentQiscusRepo: agentQiscusRepo,
	}
}

func (u *agentUsecase) CreateAgent(ctx context.Context, agent *entity.Agent) error {
	return u.agentRepo.Create(ctx, agent)
}

func (u *agentUsecase) GetAgent(ctx context.Context, agentID string) (*entity.Agent, error) {
	return u.agentRepo.GetByID(ctx, agentID)
}

func (u *agentUsecase) GetAllAgents(ctx context.Context) ([]*entity.Agent, error) {
	return u.agentRepo.GetAll(ctx)
}

func (u *agentUsecase) UpdateAgentStatus(ctx context.Context, agentID string, update *entity.AgentStatusUpdate) (*entity.Agent, error) {
	agent, err := u.agentRepo.GetByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	err = u.agentRepo.Update(ctx, agent)
	if err != nil {
		return nil, err
	}

	return agent, nil
}
