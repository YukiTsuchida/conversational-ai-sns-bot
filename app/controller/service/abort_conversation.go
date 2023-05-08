package service

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

type AbortConversationService struct {
	sns              sns.SNS
	conversationRepo conversation.ConversationRepository
}

func (svc *AbortConversationService) AbortConversation(ctx context.Context, conversationID string, reason string) error {
	err := svc.conversationRepo.Abort(ctx, conversationID, reason)
	if err != nil {
		return err
	}
	account, err := svc.sns.FetchAccountByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}
	err = svc.sns.RemoveAccountConversationID(ctx, account.ID())
	if err != nil {
		return err
	}
	return nil
}

func NewAbortConversationService(sns sns.SNS, conversationRepo conversation.ConversationRepository) *AbortConversationService {
	return &AbortConversationService{sns, conversationRepo}
}
