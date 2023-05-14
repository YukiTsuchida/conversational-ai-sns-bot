package usecases

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/services/sns"
)

type RegisterAccount struct {
	snsSvc sns.Service
}

func (uc *RegisterAccount) Execute(ctx context.Context, accountID string, credential string) error {
	err := uc.snsSvc.CreateAccount(ctx, accountID, credential)
	if err != nil {
		return err
	}
	return nil
}

func NewRegisterAccount(snsSvc sns.Service) *RegisterAccount {
	return &RegisterAccount{snsSvc}
}
