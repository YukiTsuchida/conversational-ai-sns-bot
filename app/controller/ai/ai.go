package ai

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/ai"
)

type AI interface {
	SendRequest(ctx context.Context, conversationId string) error
	SaveMessageLog(ctx context.Context, message string, role ai.MessageRole) error
}
