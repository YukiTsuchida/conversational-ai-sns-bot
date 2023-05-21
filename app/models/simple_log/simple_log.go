package simple_log

import (
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
)

type Sort string

const (
	SortAsc  Sort = "asc"
	SortDesc Sort = "desc"
)

type SimpleLog struct {
	PageIndex    int
	Size         int
	Sort         Sort
	Timezone     *time.Location
	Pages        []int
	Conversation *conversation.Conversation
	Logs         []*conversation.ConversationLog
}

func NewSimpleLog(pageIndex int, size int, sort Sort, timezone *time.Location, pages []int, conversation *conversation.Conversation, logs []*conversation.ConversationLog) *SimpleLog {
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

func (sl *SimpleLog) TimezoneStr() string {
	return sl.Timezone.String()
}
