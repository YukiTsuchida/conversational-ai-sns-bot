package service

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/sns"
)

type RegisterAccountService struct {
	sns sns.SNS
}

func (svc *RegisterAccountService) RegisterAccount(ctx context.Context, accountID string, credential string) error {
	err := svc.sns.CreateAccount(ctx, accountID, credential)
	if err != nil {
		return err
	}
	return nil
}

func NewRegisterAccountService(sns sns.SNS) *RegisterAccountService {
	return &RegisterAccountService{sns}
}
