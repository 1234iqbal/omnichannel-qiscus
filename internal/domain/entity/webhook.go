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

// Webhook payload structure based on actual Qiscus webhook
type QiscusWebhookPayload struct {
	AppID          string          `json:"app_id"`
	Source         string          `json:"source"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	AvatarURL      string          `json:"avatar_url"`
	Extras         string          `json:"extras"`
	IsResolved     bool            `json:"is_resolved"`
	RoomID         string          `json:"room_id"`
	CandidateAgent *CandidateAgent `json:"candidate_agent,omitempty"`
}

// Convert Qiscus webhook to internal RoomLog format
func (q *QiscusWebhookPayload) ToRoomLog() RoomLog {
	return RoomLog{
		RoomID:        q.RoomID,
		UserID:        q.Email,
		Name:          q.Name,
		Source:        q.Source,
		IsWaiting:     !q.IsResolved,
		Resolved:      q.IsResolved,
		UserAvatarURL: q.AvatarURL,
		Extras:        &q.Extras,
		HasNoMessage:  false,
		RoomBadge:     nil,
	}
}

// Original entities remain the same
type WebhookData struct {
	Data   ChatData `json:"data"`
	Status int      `json:"status"`
}

type ChatData struct {
	Agent        *Agent  `json:"agent"`
	AgentService *string `json:"agent_service"`
	RoomLog      RoomLog `json:"room_log"`
}

type RoomLog struct {
	ChannelID             int     `json:"channel_id"`
	Extras                *string `json:"extras"`
	HasNoMessage          bool    `json:"has_no_message"`
	IsWaiting             bool    `json:"is_waiting"`
	Name                  string  `json:"name"`
	Resolved              bool    `json:"resolved"`
	ResolvedTS            *string `json:"resolved_ts"`
	RoomBadge             *string `json:"room_badge"`
	RoomID                string  `json:"room_id"`
	Source                string  `json:"source"`
	StartServiceCommentID string  `json:"start_service_comment_id"`
	UserAvatarURL         string  `json:"user_avatar_url"`
	UserID                string  `json:"user_id"`
}

type AssignmentResponse struct {
	Status         string     `json:"status"`
	Agent          *Agent     `json:"agent,omitempty"`
	RoomID         string     `json:"room_id"`
	Message        string     `json:"message,omitempty"`
	NextAssignment *QueueItem `json:"next_assignment,omitempty"`
}
