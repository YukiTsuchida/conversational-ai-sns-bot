package usecases

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/sns"
)

type RegisterAccount struct {
	sns sns.SNS
}

func (uc *RegisterAccount) Execute(ctx context.Context, accountID string, credential string) error {
	err := uc.sns.CreateAccount(ctx, accountID, credential)
	if err != nil {
		return err
	}
	return nil
}

func NewRegisterAccount(sns sns.SNS) *RegisterAccount {
	return &RegisterAccount{sns}
}
