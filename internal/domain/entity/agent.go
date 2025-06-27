package entity

type Agent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsAvailable bool   `json:"is_available"`
}

type QiscusAgent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsAvailable bool   `json:"is_available"`
}

type GetAgentsResponse struct {
	Status int `json:"status"`
	Data   struct {
		Agents []QiscusAgent `json:"agents"`
	} `json:"data"`
}

type AssignAgentRequest struct {
	RoomID  string `json:"room_id"`
	AgentID string `json:"agent_id"`
}

type AssignAgentResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}
