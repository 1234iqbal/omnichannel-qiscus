package redis

import (
	"context"
	"encoding/json"

	"qiscus-agent-allocation/internal/domain/entity"

	"github.com/go-redis/redis/v8"
)

const (
	QueueKey = "chat_queue"
)

type QueueRepository interface {
	Enqueue(ctx context.Context, item *entity.QueueItem) error
	Dequeue(ctx context.Context) (*entity.QueueItem, error)
	GetAll(ctx context.Context) ([]*entity.QueueItem, error)
	Size(ctx context.Context) (int, error)
	Clear(ctx context.Context) error
}

type queueRepository struct {
	client *redis.Client
}

func NewQueueRepository(client *redis.Client) QueueRepository {
	return &queueRepository{
		client: client,
	}
}

func (r *queueRepository) Enqueue(ctx context.Context, item *entity.QueueItem) error {
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// LPUSH for FIFO (Left Push, Right Pop)
	return r.client.LPush(ctx, QueueKey, itemJSON).Err()
}

func (r *queueRepository) Dequeue(ctx context.Context) (*entity.QueueItem, error) {
	itemJSON, err := r.client.RPop(ctx, QueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Queue is empty
		}
		return nil, err
	}

	var item entity.QueueItem
	err = json.Unmarshal([]byte(itemJSON), &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *queueRepository) GetAll(ctx context.Context) ([]*entity.QueueItem, error) {
	itemsJSON, err := r.client.LRange(ctx, QueueKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var items []*entity.QueueItem
	for _, itemJSON := range itemsJSON {
		var item entity.QueueItem
		if err := json.Unmarshal([]byte(itemJSON), &item); err != nil {
			continue
		}
		items = append(items, &item)
	}

	return items, nil
}

func (r *queueRepository) Size(ctx context.Context) (int, error) {
	size, err := r.client.LLen(ctx, QueueKey).Result()
	return int(size), err
}

func (r *queueRepository) Clear(ctx context.Context) error {
	return r.client.Del(ctx, QueueKey).Err()
}
