package chatgpt_3_5_turbo

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
)

var _ ai.AI = (*aiChatGPT3_0TurboImpl)(nil)

type aiChatGPT3_0TurboImpl struct {
}

func (ai *aiChatGPT3_0TurboImpl) SendRequest(ctx context.Context, conversationID string) error {
	return nil
}

func (ai *aiChatGPT3_0TurboImpl) AppendSystemMessage(ctx context.Context, conversationID string, message string) error {
	return nil
}

func (ai *aiChatGPT3_0TurboImpl) AppendUserMessage(ctx context.Context, conversationID string, message string) error {
	return nil
}

func (ai *aiChatGPT3_0TurboImpl) AppendAIMessage(ctx context.Context, conversationID string, message string, purpose string) error {
	return nil
}

func NewAIChatGPT3_0TurboImpl() ai.AI {
	return &aiChatGPT3_0TurboImpl{}
}
