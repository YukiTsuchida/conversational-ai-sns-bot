package queue

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
)

type Service interface {
	Enqueue(ctx context.Context, conversationID *conversation.ID) error
}
