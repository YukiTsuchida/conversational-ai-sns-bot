package ai

import (
	"context"

	ai_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/ai"
)

type Service interface {
	SendRequest(ctx context.Context, conversationID string) error // 責務デカすぎる、実際にsendする部分はAIに依存しないので切り出す必要がある
	AppendSystemMessage(ctx context.Context, conversationID string, message *ai_model.SystemMessage) error
	AppendUserMessage(ctx context.Context, conversationID string, message *ai_model.UserMessage) error
	AppendAIMessage(ctx context.Context, conversationID string, message *ai_model.AIMessage, purpose string) error
}
