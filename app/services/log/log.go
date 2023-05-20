package log

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
)

type Service interface {
	CountMessageLog(ctx context.Context, conversationID *conversation.ID) (int, error)
	FetchMessageLogs(ctx context.Context, conversationID *conversation.ID, page int, size int, sort string) ([]*conversation.ConversationLog, error)
}
