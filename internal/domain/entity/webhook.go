package entity

// Webhook payload structure based on actual Qiscus webhook
type QiscusWebhookPayload struct {
	AppID          string          `json:"app_id"`
	Source         string          `json:"source"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	AvatarURL      string          `json:"avatar_url"`
	Extras         string          `json:"extras"`
	IsResolved     bool            `json:"is_resolved"`
	LatestService  LatestService   `json:"latest_service"`
	RoomID         string          `json:"room_id"`
	CandidateAgent *CandidateAgent `json:"candidate_agent,omitempty"`
}

type LatestService struct {
	ID                    int     `json:"id"`
	UserID                int     `json:"user_id"`
	RoomLogID             int     `json:"room_log_id"`
	AppID                 int     `json:"app_id"`
	RoomID                string  `json:"room_id"`
	Notes                 *string `json:"notes"`
	ResolvedAt            *string `json:"resolved_at"`
	IsResolved            bool    `json:"is_resolved"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
	FirstCommentID        string  `json:"first_comment_id"`
	LastCommentID         string  `json:"last_comment_id"`
	RetrievedAt           string  `json:"retrieved_at"`
	FirstCommentTimestamp *string `json:"first_comment_timestamp"`
}

type CandidateAgent struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Email               string   `json:"email"`
	AuthenticationToken string   `json:"authentication_token"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	SDKEmail            string   `json:"sdk_email"`
	SDKKey              string   `json:"sdk_key"`
	IsAvailable         bool     `json:"is_available"`
	Type                int      `json:"type"`
	AvatarURL           string   `json:"avatar_url"`
	AppID               int      `json:"app_id"`
	IsVerified          bool     `json:"is_verified"`
	NotificationsRoomID string   `json:"notifications_room_id"`
	BubbleColor         string   `json:"bubble_color"`
	QismoKey            string   `json:"qismo_key"`
	DirectLoginToken    *string  `json:"direct_login_token"`
	TypeAsString        string   `json:"type_as_string"`
	AssignedRules       []string `json:"assigned_rules"`
}

// Convert Qiscus webhook to internal RoomLog format
func (q *QiscusWebhookPayload) ToRoomLog() RoomLog {
	return RoomLog{
		RoomID:                q.RoomID,
		UserID:                q.Email,
		Name:                  q.Name,
		Source:                q.Source,
		IsWaiting:             !q.IsResolved,
		Resolved:              q.IsResolved,
		UserAvatarURL:         q.AvatarURL,
		Extras:                &q.Extras,
		HasNoMessage:          false,
		RoomBadge:             nil,
		ResolvedTS:            q.LatestService.ResolvedAt,
		StartServiceCommentID: q.LatestService.FirstCommentID,
		ChannelID:             q.LatestService.AppID,
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
