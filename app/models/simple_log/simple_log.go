package simple_log

import "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

type SimpleLog struct {
	ConversationID string
	Pages          []int
	Logs           []*conversation.ConversationLog
}

func NewSimpleLog(id conversation.ID, pages []int, logs []*conversation.ConversationLog) *SimpleLog {
	return &SimpleLog{
		ConversationID: id.ToString(),
		Pages:          pages,
		Logs:           logs,
	}
}
