package usecases

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/models/conversation"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
)

type AbortConversation struct {
	snsSvc           sns.Service
	conversationRepo repositories.Conversation
}

func (uc *AbortConversation) Execute(ctx context.Context, conversationID *conversation.ID, reason string) error {
	err := uc.conversationRepo.Abort(ctx, conversationID, reason)
	if err != nil {
		return err
	}
	account, err := uc.snsSvc.FetchAccountByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}
	err = uc.snsSvc.RemoveAccountConversationID(ctx, account.ID())
	if err != nil {
		return err
	}
	return nil
}

func NewAbortConversation(snsSvc sns.Service, conversationRepo repositories.Conversation) *AbortConversation {
	return &AbortConversation{snsSvc, conversationRepo}
}
