package entity

type QiscusWebhook struct {
	AppID          string          `json:"app_id"`
	RoomID         string          `json:"room_id"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Source         string          `json:"source"`
	IsResolved     bool            `json:"is_resolved"`
	CandidateAgent *CandidateAgent `json:"candidate_agent,omitempty"`
}

type CandidateAgent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsAvailable bool   `json:"is_available"`
}

type QiscusResolvedWebhook struct {
	AppCode  string `json:"app_code"`
	Customer struct {
		Name   string `json:"name"`
		UserID string `json:"user_id"`
		Avatar string `json:"avatar"`
	} `json:"customer"`
	ResolvedBy struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		IsAvailable bool   `json:"is_available"`
		Type        string `json:"type"`
	} `json:"resolved_by"`
	Service struct {
		ID         int    `json:"id"`
		RoomID     string `json:"room_id"`
		IsResolved bool   `json:"is_resolved"`
		Source     string `json:"source"`
		Notes      string `json:"notes"`
	} `json:"service"`
}
