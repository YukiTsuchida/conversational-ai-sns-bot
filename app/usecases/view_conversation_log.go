package usecases

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/simple_log"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/log"
)

type ViewConversation struct {
	logSvc           log.Service
	conversationRepo repositories.Conversation
}

func (uc *ViewConversation) Execute(ctx context.Context, conversationId *conversation.ID, page int, size int, sort string, timezone string) (*simple_log.SimpleLog, error) {
	conversation, err := uc.conversationRepo.FetchByID(ctx, conversationId)
	if err != nil {
		return nil, err
	}

	logCount, err := uc.logSvc.CountMessageLog(ctx, conversationId)
	if err != nil {
		return nil, err
	}

	logs, err := uc.logSvc.FetchMessageLogs(ctx, conversationId, page, size, sort)
	if err != nil {
		return nil, err
	}

	return simple_log.NewSimpleLog(
		conversation.ID,
		logCount/size,
		logs,
	), nil
}

func NewViewConversationLog(logSvc log.Service, conversationRepo repositories.Conversation) *ViewConversation {
	return &ViewConversation{logSvc, conversationRepo}
}
