package ai

import (
	"context"
)

type AI interface {
	SendRequest(ctx context.Context, conversationID string) error
	AppendSystemMessage(ctx context.Context, conversationID string, message string) error
	AppendUserMessage(ctx context.Context, conversationID string, message string) error
	AppendAIMessage(ctx context.Context, conversationID string, message string, purpose string) error
}
