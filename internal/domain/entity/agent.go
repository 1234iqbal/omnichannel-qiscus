package entity

import "time"

type Agent struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	IsOnline      bool      `json:"is_online"`
	MaxConcurrent int       `json:"max_concurrent"`
	CurrentChats  int       `json:"current_chats"`
	LastAssigned  time.Time `json:"last_assigned"`
}

type AgentStatusUpdate struct {
	IsOnline      *bool `json:"is_online"`
	MaxConcurrent *int  `json:"max_concurrent"`
}

func (a *Agent) CanTakeMoreChats() bool {
	return a.IsOnline && a.CurrentChats < a.MaxConcurrent
}

func (a *Agent) AssignChat() {
	a.CurrentChats++
	a.LastAssigned = time.Now()
}

func (a *Agent) ResolveChat() {
	if a.CurrentChats > 0 {
		a.CurrentChats--
	}
}
