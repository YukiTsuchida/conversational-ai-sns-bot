package usecases

import (
	"context"
	"time"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/simple_log"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/ai"
)

type ViewConversation struct {
	aiSvc            ai.Service
	conversationRepo repositories.Conversation
}

func (uc *ViewConversation) Execute(ctx context.Context, conversationId *conversation.ID, page int, size int, sort conversation.Sort, timezone *time.Location) (*simple_log.SimpleLog, error) {
	conversation, err := uc.conversationRepo.FetchByID(ctx, conversationId)
	if err != nil {
		// TODO: 検索結果がない場合もエラーを返してしまう
		return nil, err
	}

	logCount, err := uc.aiSvc.CountMessageLog(ctx, conversationId)
	if err != nil {
		return nil, err
	}

	logs, err := uc.aiSvc.FetchMessageLogs(ctx, conversationId, page, size, sort)
	if err != nil {
		return nil, err
	}
	for _, log := range logs {
		log.SetTimezone(timezone)
	}

	var pages []int
	for i := 0; i <= logCount/size; i++ {
		pages = append(pages, i)
	}

	return simple_log.NewSimpleLog(
		page,
		size,
		sort,
		timezone,
		pages,
		conversation,
		logs,
	), nil
}

func NewViewConversationLog(logSvc ai.Service, conversationRepo repositories.Conversation) *ViewConversation {
	return &ViewConversation{logSvc, conversationRepo}
}
