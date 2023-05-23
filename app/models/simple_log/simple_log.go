package simple_log

import (
	"bytes"
	"text/template"
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
)

const templateDir = "/app/http/template/simple_conversation_log_viewer.html"

type SimpleLog struct {
	PageIndex    int
	Size         int
	Sort         conversation.Sort
	Timezone     *time.Location
	Pages        []int
	Conversation *conversation.Conversation
	Logs         []*conversation.ConversationLog
}

func NewSimpleLog(pageIndex int, size int, sort conversation.Sort, timezone *time.Location, pages []int, conversation *conversation.Conversation, logs []*conversation.ConversationLog) *SimpleLog {
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

func (sl *SimpleLog) GenerateHtml() ([]byte, error) {
	writer := new(bytes.Buffer)

	t, err := template.ParseFiles(templateDir)
	if err != nil {
		return nil, err
	}

	err = t.Execute(writer, sl)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
