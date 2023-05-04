package chatgpt_3_5_turbo

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ai"
	ai_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/ai"
)

var _ ai.AI = (*aiChatGPT3_0TurboImpl)(nil)

type aiChatGPT3_0TurboImpl struct {
}

func (ai *aiChatGPT3_0TurboImpl) SendRequest(ctx context.Context, conversationId string) error {
	return nil
}

func (ai *aiChatGPT3_0TurboImpl) SaveMessageLog(ctx context.Context, message string, role ai_model.MessageRole) error {
	return nil
}

func NewAIChatGPT3_0TurboImpl() ai.AI {
	return &aiChatGPT3_0TurboImpl{}
}
