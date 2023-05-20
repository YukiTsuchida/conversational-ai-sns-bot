package conversation

import "time"

type LogID struct {
	id string
}

func NewLogID(id string) *LogID {
	return &LogID{id: id}
}

func (id *LogID) ToString() string {
	return id.id
}

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type ConversationLog struct {
	LogID
	message   string
	purpose   string
	role      Role
	createdAt time.Time
}

func NewConversationLog(logId string, message string, purpose string, role string, createdAt time.Time) *ConversationLog {
	return &ConversationLog{
		LogID:     LogID{id: logId},
		message:   message,
		purpose:   purpose,
		role:      Role(role),
		createdAt: createdAt,
	}
}
