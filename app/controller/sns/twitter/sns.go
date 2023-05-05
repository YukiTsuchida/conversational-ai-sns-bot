package twitter

import (
	"context"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/cmd"

	sns_model "github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/model/sns"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/sns"
)

var _ sns.SNS = (*snsTwitterImpl)(nil)

type snsTwitterImpl struct {
}

func (sns *snsTwitterImpl) GetAccountById(ctx context.Context, accountId string) (*sns_model.Account, error) {
	return nil, nil
}

func (sns *snsTwitterImpl) CreateAccount(ctx context.Context, accountId string, credential string) error {
	return nil
}

func (sns *snsTwitterImpl) GiveAccountConversationId(ctx context.Context, conversationId string) error {
	return nil
}

func (sns *snsTwitterImpl) ExecuteCmd(ctx context.Context, cmd *cmd.Command) (*sns_model.Response, error) {
	return nil, nil
}

func NewSNSTwitterImpl() sns.SNS {
	return &snsTwitterImpl{}
}
