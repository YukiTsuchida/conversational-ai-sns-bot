package conversation

import (
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
)

type LogID struct {
	id string
}

func NewLogID(id string) *LogID {
	return &LogID{id: id}
}

func (id *LogID) ToString() string {
	return id.id
}

type ConversationLog struct {
	LogID
	message   ai.Message
	purpose   string
	role      ai.Role
	createdAt *time.Time
	timezone  *time.Location
}

func NewConversationLog(logId string, message ai.Message, purpose string, role ai.Role, createdAt *time.Time) *ConversationLog {
	return &ConversationLog{
		LogID:     LogID{id: logId},
		message:   message,
		purpose:   purpose,
		role:      role,
		createdAt: createdAt,
	}
}

func (c *ConversationLog) LogIDStr() string {
	return c.LogID.ToString()
}

func (c *ConversationLog) MessageStr() string {
	return c.message.ToString()
}

func (c *ConversationLog) Purpose() string {
	return c.purpose
}

func (c *ConversationLog) Role() string {
	return string(c.role)
}

func (c *ConversationLog) CreatedAtStr() string {
	if c.timezone == nil {
		return c.createdAt.Format("2006-01-02 15:04:05 -07:00")
	}
	return c.createdAt.In(c.timezone).Format("2006-01-02 15:04:05 -07:00")
}

func (c *ConversationLog) SetTimezone(timezone *time.Location) {
	c.timezone = timezone
}
