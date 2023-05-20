package ai

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"

	ai_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
)

type Service interface {
	SendRequest(ctx context.Context, conversationID *conversation.ID) (*ai_model.Response, error)
	AppendSystemMessage(ctx context.Context, conversationID *conversation.ID, message *ai_model.SystemMessage) error
	AppendUserMessage(ctx context.Context, conversationID *conversation.ID, message *ai_model.UserMessage) error
	AppendAIMessage(ctx context.Context, conversationID *conversation.ID, message *ai_model.AIMessage, purpose string) error
	CountMessageLog(ctx context.Context, conversationID *conversation.ID) (int, error)
	FetchMessageLogs(ctx context.Context, conversationID *conversation.ID, page int, size int, sort string) ([]*conversation.ConversationLog, error)
}
