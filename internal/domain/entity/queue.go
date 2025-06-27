package entity

import "time"

type QueueItem struct {
	CustomerID string    `json:"customer_id"`
	RoomID     string    `json:"room_id"`
	Channel    string    `json:"channel"`
	Timestamp  time.Time `json:"timestamp"`
}
