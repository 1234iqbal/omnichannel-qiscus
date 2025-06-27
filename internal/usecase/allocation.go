package usecase

import (
	"context"
	"fmt"
	"log"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/internal/repository/qiscus"
	"qiscus-agent-allocation/internal/repository/redis"
)

type AllocationUsecase interface {
	ProcessIncomingMessage(ctx context.Context, roomLog entity.RoomLog) (*entity.AssignmentResponse, error)
	ProcessResolvedMessage(ctx context.Context, roomLog entity.RoomLog) (*entity.AssignmentResponse, error)
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

func (u *allocationUsecase) ProcessIncomingMessage(ctx context.Context, roomLog entity.RoomLog) (*entity.AssignmentResponse, error) {

	// Step 1: Check if room already has an assigned agent
	existingAssignment, err := u.agentRepo.FindAgentByRoomID(ctx, roomLog.RoomID)
	if err == nil && existingAssignment != "" {
		// Get full agent details using the agent ID
		existingAgent, err := u.agentRepo.GetByID(ctx, existingAssignment)
		if err == nil {
			log.Printf("Room %s already assigned to agent %s", roomLog.RoomID, existingAgent.Email)
			return &entity.AssignmentResponse{
				Status:  "existing_assignment",
				RoomID:  roomLog.RoomID,
				Message: fmt.Sprintf("Room already assigned to agent %s", existingAgent.Name),
				Agent:   existingAgent, // Now this is *entity.Agent
			}, nil
		}
	}

	// // Step 2: Get available agents from Qiscus API
	availableAgents, err := u.agentQiscusRepo.GetAvailableAgents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available agents: %w", err)
	}

	if len(availableAgents) == 0 {
		log.Printf("No available agents found")
		return &entity.AssignmentResponse{
			Status:  "no_agents_available",
			RoomID:  roomLog.RoomID,
			Message: "No agents are currently available",
		}, nil
	}

	// // Step 3: Get all active agents from Redis
	allLocalAgents, err := u.agentRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get local agents: %w", err)
	}

	if len(allLocalAgents) == 0 {
		log.Printf("No agents found in local storage")
		return &entity.AssignmentResponse{
			Status:  "no_agents_configured",
			RoomID:  roomLog.RoomID,
			Message: "No agents are configured in the system",
		}, nil
	}

	// // Step 4: Filter available agents (online and not at max capacity)

	var availableAgentsCapacity []*entity.Agent
	// Hitung jumlah room yang ditangani oleh setiap agent
	roomCount := make(map[string]int)
	for _, a := range assignments {
		key := fmt.Sprintf("%d:%s", a.CompanyID, a.AgentName)
		roomCount[key]++
	}

	// Tampilkan agent yang menangani kurang dari 2 room
	fmt.Println("Agents handling fewer than 2 rooms:")
	for _, ag := range allLocalAgents {
		key := fmt.Sprintf("%d:%s", ag.CompanyID, ag.Name)
		if roomCount[key] < 2 {
			fmt.Printf("- %s (rooms: %d)\n", ag.Name, roomCount[key])
		}
	}

	println("=================")
	println("step 4")
	println(availableAgentsCapacity)
	println("=================")

	// if len(availableAgents) == 0 {
	// 	log.Printf("No available agents found, adding to queue")
	// 	// Step 6: If no available agents, add to queue
	// 	queueItem := entity.NewQueueItem(roomLog)
	// 	err = u.queueRepo.Enqueue(ctx, queueItem)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to queue message: %w", err)
	// 	}

	// 	log.Printf("Added room %s to queue", roomLog.RoomID)
	// 	return &entity.AssignmentResponse{
	// 		Status:  "queued",
	// 		RoomID:  roomLog.RoomID,
	// 		Message: "No agents available, added to queue",
	// 	}, nil
	// }

	// // Step 5: Select best available agent (load balancing)
	// selectedAgent := u.selectBestAgent(availableAgents)
	// if selectedAgent == nil {
	// 	return nil, fmt.Errorf("failed to select agent")
	// }

	// // Assign chat to agent locally first
	// selectedAgent.AssignChat()
	// err = u.agentRepo.Update(ctx, selectedAgent)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to update agent locally: %w", err)
	// }

	// // Add room to agent's chat list
	// err = u.agentRepo.AddChatToAgent(ctx, selectedAgent.ID, roomLog.RoomID)
	// if err != nil {
	// 	// Rollback agent update
	// 	selectedAgent.ResolveChat()
	// 	u.agentRepo.Update(ctx, selectedAgent)
	// 	return nil, fmt.Errorf("failed to add chat to agent: %w", err)
	// }

	// // Assign agent to room via Qiscus API
	// err = u.qiscusService.AssignAgentToRoom(ctx, roomLog.RoomID, selectedAgent.ID)
	// if err != nil {
	// 	// Rollback local assignment if Qiscus API fails
	// 	log.Printf("Failed to assign agent in Qiscus, rolling back: %v", err)
	// 	u.rollbackAgentAssignment(ctx, selectedAgent.ID, roomLog.RoomID)

	// 	// Add to queue as fallback
	// 	queueItem := entity.NewQueueItem(roomLog)
	// 	u.queueRepo.Enqueue(ctx, queueItem)

	// 	return &entity.AssignmentResponse{
	// 		Status:  "assignment_failed_queued",
	// 		RoomID:  roomLog.RoomID,
	// 		Message: "Failed to assign via Qiscus API, added to queue",
	// 	}, nil
	// }

	log.Printf("Successfully assigned room %s to agent %s (%s)",
		roomLog.RoomID, roomLog.Name)

	return &entity.AssignmentResponse{
		Status:  "assigned",
		RoomID:  roomLog.RoomID,
		Message: fmt.Sprintf("Assigned to agent %s", roomLog.Name),
	}, nil
}

func (u *allocationUsecase) ProcessResolvedMessage(ctx context.Context, roomLog entity.RoomLog) (*entity.AssignmentResponse, error) {
	// Find which agent was handling this chat
	// // agentID, err := u.agentRepo.FindAgentByRoomID(ctx, roomLog.RoomID)
	// // if err != nil {
	// // 	return nil, fmt.Errorf("agent not found for room %s: %w", roomLog.RoomID, err)
	// // }

	// // // Remove chat from agent
	// // err = u.removeChatFromAgent(ctx, agentID, roomLog.RoomID)
	// // if err != nil {
	// // 	log.Printf("Error removing chat from agent: %v", err)
	// // }

	response := &entity.AssignmentResponse{
		Status: "resolved",
		RoomID: roomLog.RoomID,
	}

	// // // Try to assign next chat from queue to this agent
	// // nextChat, err := u.assignNextFromQueue(ctx, agentID)
	// if err != nil {
	// 	log.Printf("No chats in queue or error: %v", err)
	// } else if nextChat != nil {
	// 	log.Printf("Assigned next chat %s to agent %s", nextChat.RoomID, agentID)
	// 	response.NextAssignment = nextChat
	// }

	return response, nil
}
