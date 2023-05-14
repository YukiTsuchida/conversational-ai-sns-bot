package ai

import (
	"context"
)

type Service interface {
	SendRequest(ctx context.Context, conversationID string) error // 責務デカすぎる、実際にsendする部分はAIに依存しないので切り出す必要がある
	AppendSystemMessage(ctx context.Context, conversationID string, message string) error
	AppendUserMessage(ctx context.Context, conversationID string, message string) error
	AppendAIMessage(ctx context.Context, conversationID string, message string, purpose string) error
}
