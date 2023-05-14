package usecases

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/repositories"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/sns"
)

type AbortConversation struct {
	sns              sns.SNS
	conversationRepo repositories.Conversation
}

func (uc *AbortConversation) Execute(ctx context.Context, conversationID string, reason string) error {
	err := uc.conversationRepo.Abort(ctx, conversationID, reason)
	if err != nil {
		return err
	}
	account, err := uc.sns.FetchAccountByConversationID(ctx, conversationID)
	if err != nil {
		return err
	}
	err = uc.sns.RemoveAccountConversationID(ctx, account.ID())
	if err != nil {
		return err
	}
	return nil
}

func NewAbortConversation(sns sns.SNS, conversationRepo repositories.Conversation) *AbortConversation {
	return &AbortConversation{sns, conversationRepo}
}
