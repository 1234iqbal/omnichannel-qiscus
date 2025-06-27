package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"qiscus-agent-allocation/internal/domain/entity"
	"qiscus-agent-allocation/internal/usecase"
)

type WorkerService struct {
	allocationUsecase usecase.AllocationUsecase
}

func NewWorkerService(allocationUsecase usecase.AllocationUsecase) *WorkerService {
	return &WorkerService{
		allocationUsecase: allocationUsecase,
	}
}

func (w *WorkerService) Start(ctx context.Context) {
	log.Println("Worker service started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker service stopped")
			return
		default:
			w.processQueue()
		}
	}
}

func (w *WorkerService) processQueue() {
	// 1. Check Redis Queue (RPOP)
	queueData, err := w.allocationUsecase.GetFromQueue()
	if err != nil || queueData == "" {
		// Queue empty, wait and try again
		time.Sleep(5 * time.Second)
		return
	}

	log.Printf("Processing queue item: %s", queueData)

	// 2. Extract customer request
	var item entity.QueueItem
	if err := json.Unmarshal([]byte(queueData), &item); err != nil {
		log.Printf("Failed to parse queue item: %v", err)
		return
	}

	// 3. Fetch online agents from Qiscus API
	agents, err := w.allocationUsecase.GetOnlineAgents()
	if err != nil {
		log.Printf("Failed to get online agents: %v", err)
		// Return to queue
		w.allocationUsecase.AddToQueue(item)
		time.Sleep(5 * time.Second)
		return
	}

	if len(agents) == 0 {
		log.Println("No online agents available")
		// Return to queue
		w.allocationUsecase.AddToQueue(item)
		time.Sleep(5 * time.Second)
		return
	}

	// 4. Check agent capacity and filter available agents
	availableAgent := w.findAvailableAgent(agents)
	if availableAgent == nil {
		log.Println("No available agents (all at capacity)")
		// Return to queue
		w.allocationUsecase.AddToQueue(item)
		time.Sleep(5 * time.Second)
		return
	}

	// 5. Assign agent via Qiscus API
	err = w.allocationUsecase.AssignAgent(item.RoomID, availableAgent.ID)
	if err != nil {
		log.Printf("Failed to assign agent: %v", err)
		// Return to queue
		w.allocationUsecase.AddToQueue(item)
		time.Sleep(5 * time.Second)
		return
	}

	// 6. Update agent capacity
	err = w.allocationUsecase.IncrementAgentCapacity(availableAgent.ID)
	if err != nil {
		log.Printf("Failed to update agent capacity: %v", err)
	}

	// 7. Log successful assignment
	log.Printf("Successfully assigned agent %s to customer %s (room: %s)",
		availableAgent.ID, item.CustomerID, item.RoomID)
}

func (w *WorkerService) findAvailableAgent(agents []entity.Agent) *entity.Agent {
	maxCapacity := 2 // Default max customers per agent

	for _, agent := range agents {
		// Check current capacity from Redis
		currentCapacity, err := w.allocationUsecase.GetAgentCapacity(agent.ID)
		if err != nil {
			log.Printf("Failed to get capacity for agent %s: %v", agent.ID, err)
			continue
		}

		// If agent has available slots
		if currentCapacity < maxCapacity {
			return &agent
		}
	}

	return nil
}
