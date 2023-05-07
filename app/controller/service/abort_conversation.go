package service

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/conversation"
)

type AbortConversationService struct {
	conversationRepo conversation.ConversationRepository
}

func (svc *AbortConversationService) AbortConversation(ctx context.Context, conversationID string, reason string) error {
	err := svc.conversationRepo.Abort(ctx, conversationID, reason)
	if err != nil {
		return err
	}
	return nil
}

func NewAbortConversationService(conversationRepo conversation.ConversationRepository) *AbortConversationService {
	return &AbortConversationService{conversationRepo}
}
