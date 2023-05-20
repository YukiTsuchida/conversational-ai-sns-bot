package simple_log

import "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

type SimpleLog struct {
	ConversationID string
	Page           int
	Logs           []*conversation.ConversationLog
}

func NewSimpleLog(id conversation.ID, page int, logs []*conversation.ConversationLog) *SimpleLog {
	return &SimpleLog{
		ConversationID: id.ToString(),
		Page:           page,
		Logs:           logs,
	}
}
