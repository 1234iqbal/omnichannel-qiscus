package entity

import "time"

type QueueItem struct {
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

func NewQueueItem(roomLog RoomLog) *QueueItem {
	return &QueueItem{
		RoomID:    roomLog.RoomID,
		UserID:    roomLog.UserID,
		Name:      roomLog.Name,
		Source:    roomLog.Source,
		Timestamp: time.Now(),
	}
}
