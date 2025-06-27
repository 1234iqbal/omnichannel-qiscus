package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	QueueKey = "chat_queue"
)

type QueueRepository interface {
	Push(data string) error
	Pop() (string, error)
	Exists(roomID, channel, customerID string) (bool, error)
}

type queueRepository struct {
	client *redis.Client
}

func NewQueueRepository(client *redis.Client) QueueRepository {
	return &queueRepository{
		client: client,
	}
}

func (r *queueRepository) Push(data string) error {
	ctx := context.Background()

	// LPUSH adds to the left (beginning) of the list
	err := r.client.LPush(ctx, QueueKey, data).Err()
	if err != nil {
		return fmt.Errorf("failed to push to queue: %w", err)
	}

	return nil
}

func (r *queueRepository) Pop() (string, error) {
	ctx := context.Background()

	// RPOP removes and returns element from the right (end) of the list
	// This gives us FIFO behavior when combined with LPUSH
	result := r.client.RPop(ctx, QueueKey)

	// Check if queue is empty
	if result.Err() == redis.Nil {
		return "", fmt.Errorf("queue is empty")
	}

	if result.Err() != nil {
		return "", fmt.Errorf("failed to pop from queue: %w", result.Err())
	}

	return result.Val(), nil
}

// Exists checks if room_id already exists in queue
func (r *queueRepository) Exists(roomID, channel, customerID string) (bool, error) {
	ctx := context.Background()

	// Get all items in queue without removing them (LRANGE)
	result := r.client.LRange(ctx, QueueKey, 0, -1)
	if result.Err() != nil {
		return false, fmt.Errorf("failed to check queue: %w", result.Err())
	}

	items := result.Val()

	// Check each item in queue
	for _, item := range items {
		var queueItem map[string]interface{}
		json.Unmarshal([]byte(item), &queueItem)

		// Get fields from queue item
		queueRoomID, _ := queueItem["room_id"].(string)
		queueChannel, _ := queueItem["channel"].(string)
		queueCustomerID, _ := queueItem["customer_id"].(string)

		// Check exact match of all 3 fields
		if queueRoomID == roomID && queueChannel == channel && queueCustomerID == customerID {
			return true, nil
		}
	}

	return false, nil
}
