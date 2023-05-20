package simple_log

import "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

type SimpleLog struct {
	PageIndex    int
	Size         int
	Sort         string
	Timezone     string
	Pages        []int
	Conversation *conversation.Conversation
	Logs         []*conversation.ConversationLog
}

func NewSimpleLog(pageIndex int, size int, sort string, timezone string, pages []int, conversation *conversation.Conversation, logs []*conversation.ConversationLog) *SimpleLog {
	return &SimpleLog{
		PageIndex:    pageIndex,
		Size:         size,
		Sort:         sort,
		Timezone:     timezone,
		Pages:        pages,
		Conversation: conversation,
		Logs:         logs,
	}
}
